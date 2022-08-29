package server_tcp

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "3000"
	CONN_TYPE = "tcp"
)

// Info Metadata
type Info struct {
	CurrFile string
	Filesize int64
	Total    int64
	Clients  map[int]string
	MaxSize  int64
}

var chans = make([]chan []byte, 3)

var Data = make(map[int]Info)

var m sync.Mutex

var MAX_SIZE int64

func Run(numChannels int, maxFilesize int64) {
	fmt.Printf("Server starting...Channels: %d Max file size: %d(B)\n", numChannels, maxFilesize)

	// Create and initialize every channel and the Data struct.
	for i := 0; i < numChannels; i++ {
		chans = append(chans, make(chan []byte))
		Data[i] = Info{
			"",
			0,
			0,
			make(map[int]string),
			maxFilesize,
		}
	}
	MAX_SIZE = maxFilesize
	fmt.Println("Server running...2")

	listen, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer listen.Close()
	fmt.Printf("Waiting for incoming requests on %s:%s\n", CONN_HOST, CONN_PORT)
	go startAPI()
	// forever...
	for i := 0; ; i++ {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(2)
		}
		go handleIncomingRequest(conn, i)
	}
}

func handleIncomingRequest(conn net.Conn, id int) {
	// store incoming data
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		return
	}

	processRequest(string(buffer[:n]), conn, id)
	conn.Close()
}

/*
Process any given request sent to the server.
A valid request has the form:

	Sending files: -> <content-size> <file> <channel>
	Subscribing to channel: listen <channel>
*/
func processRequest(body string, conn net.Conn, id int) {
	content := strings.Split(body, " ")
	fmt.Printf(">>>>>> %v\n", body)

	channel, _ := strconv.Atoi(content[len(content)-2])
	clientAddr := content[len(content)-1]

	// Spaghetti to detect when a client terminates.
	go func() {
		b := make([]byte, 1)
		for {
			_, err := conn.Read(b)
			if err != nil {

				if err == io.EOF {
					fmt.Printf("Client %d disconnected from channel %d\n", id, channel)
					m.Lock()
					delete(Data[channel].Clients, id)
					m.Unlock()
					b = nil
				}
				return
			}
		}
	}()

	switch {
	case content[0] == "->":

		addClient(id, channel, clientAddr)
		receiveFile(content[1], content[2], channel, id, conn)

		m.Lock()
		delete(Data[channel].Clients, id)
		m.Unlock()

	case content[0] == "listen":
		addClient(id, channel, clientAddr)
		sendtoClient(channel, id, conn)

	default:
		fmt.Printf("Malformed request: %s\n", body)
		os.Exit(1)
	}
}

/*
Receives files from the connected client under the specified channel.
*/
func receiveFile(size, filename string, channel, connId int, conn net.Conn) {
	// Convert size to int64
	fileSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		fmt.Println("Error reading file size")
		return
	}

	if fileSize > MAX_SIZE {
		fmt.Printf("Error filesize (%d) exceeds maximum filesize allowed (%d)\n", fileSize, MAX_SIZE)
		return
	}

	if _, err = conn.Write([]byte("OK")); err != nil {
		checkError(err, connId, channel)
	}

	m.Lock()
	if copy, ok := Data[channel]; ok {
		copy.CurrFile = filename
		copy.Filesize = fileSize
		copy.Total = Data[channel].Total + fileSize

		Data[channel] = copy
	}
	m.Unlock()
	inputBuffer := make([]byte, fileSize)
	_, err = conn.Read(inputBuffer)

	checkError(err, connId, channel)

	fmt.Printf("Emitting data over channel %d\n", channel)
	chans[channel] <- inputBuffer
}

/*
Sends a file to the clients listening on the specified channel.
*/
func sendtoClient(channel, connId int, conn net.Conn) {

	if _, ok := Data[channel]; !ok {
		fmt.Printf("Channel %d does not exist.\n", channel)
		return
	}
	fmt.Printf("Subscribing to %d\n", channel)

	// Receives and broadcast file contents to the channel.
	for {
		buf := make([]byte, 2)
		select {
		case data := <-chans[channel]:
			m.Lock()
			fmt.Printf("Sending %v\n", Data[channel].CurrFile)
			n, err := conn.Write([]byte(Data[channel].CurrFile))
			m.Unlock()
			if err != nil {
				fmt.Println("Couldn't send filename to client", err)
				return
			}
			fmt.Printf("Waiting on OK %d\n", n)
			if _, err := conn.Read(buf); err != nil {
				checkError(err, connId, channel)
				fmt.Println("Couldn't read OK from client")
				return
			}

			if string(buf) == "OK" {
				n, err := conn.Write(data)
				checkError(err, connId, channel)

				fmt.Printf("Sent %dB to connection\n", n)

				if _, err := conn.Read(buf); err != nil {
					checkError(err, connId, channel)
					fmt.Println("Couldn't read ok from client", err)
					return
				}
			}
		}
	}
}
func addClient(id, channel int, addr string) {
	m.Lock()
	if copy, ok := Data[channel]; ok {
		copy.Clients[id] = addr
		fmt.Printf("%v\n", copy.Clients)
		Data[channel] = copy
	}
	m.Unlock()
}
