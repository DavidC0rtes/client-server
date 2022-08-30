package server_tcp

import (
	"fmt"
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

var chans []chan []byte

var Data = make(map[int]Info)

var m sync.Mutex

var MAX_SIZE int64

// Run, starts the server with the desired number of channels and desired max filesize.
func Run(numChannels int, maxFilesize int64) {
	fmt.Printf("Server starting with %d channels and max file size of %d(B)\n", numChannels, maxFilesize)

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
		// i serves as an id for the connection, this will be useful
		// to know which client disconnects.
		go handleIncomingRequest(conn, i)
	}
}

func handleIncomingRequest(conn net.Conn, id int) {
	// store incoming data
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading protocol message")
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

	channel, _ := strconv.Atoi(content[len(content)-1])
	clientAddr := conn.RemoteAddr().String()

	// Spaghetti to detect when a client terminates.
	// When the client sends SIGTERM (Ctrl-C) the server
	// receives an EOF error
	/* go func() {
		b := make([]byte, 1) // We don't read anything from b, we just need to catch an error.
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
	}() */

	switch {
	case content[0] == "->": // If we are receiving from a client...

		addClient(id, channel, clientAddr)
		receiveFile(content[1], content[2], channel, id, conn)

		// Every access to Data needs to be inside a mutual exclusion block
		// bc it is a global variable and is not thread-safe
		m.Lock()
		delete(Data[channel].Clients, id)
		m.Unlock()

	case content[0] == "listen": // If a client wants to listen...
		addClient(id, channel, clientAddr)
		sendtoClient(channel, id, conn)

	default:
		fmt.Printf("Malformed request: %s\n", body)
		os.Exit(1)
	}
}

// Receives files from the connected client under the specified channel.
func receiveFile(size, filename string, channel, connId int, conn net.Conn) {
	// Convert size to int64
	fileSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		fmt.Println("Error reading filesize")
		return
	}

	if fileSize > MAX_SIZE {
		fmt.Printf("Error filesize (%d) exceeds maximum filesize allowed (%d)\n", fileSize, MAX_SIZE)
		return
	}

	if _, err = conn.Write([]byte("OK")); err != nil {
		fmt.Println("Couldn't send OK to client", err)
		return
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
	if _, err = conn.Read(inputBuffer); err != nil {
		fmt.Println("Error reading from input buffer", err)
		return
	}

	fmt.Printf("Emitting data over channel %d\n", channel)
	for {
		chans[channel] <- inputBuffer
	}
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

loop:
	for {
		buf := make([]byte, 2)
		select {
		case data := <-chans[channel]:
			fmt.Printf("Sending %v\n", Data[channel].CurrFile)
			m.Lock()
			n, err := conn.Write([]byte(Data[channel].CurrFile))
			m.Unlock()
			if err != nil {
				fmt.Println("Couldn't send filename to client", err)
				return
			}
			fmt.Printf("Waiting on OK %d\n", n)
			if _, err := conn.Read(buf); err != nil {
				fmt.Println("Couldn't read OK from client")
				return
			}

			if string(buf) == "OK" {
				n, err = conn.Write(data)
				if err != nil {
					fmt.Println("Couldn't send data to client", err)
				}

				fmt.Printf("Sent %dB to connection\n", n)

				if _, err := conn.Read(buf); err != nil {
					fmt.Println("Couldn't read response from client", err)
					return
				}
			}
		default:
			conn.Write([]byte("Nofile"))
			bf := make([]byte, 1)
			n, err := conn.Read(bf)
			if err != nil {
				fmt.Println("Couldn't read response from client", err)
			}

			resp := string(bf[:n])
			if resp != "n" {
				fmt.Printf("Client %s whishes to disconnect.\n", conn.RemoteAddr().String())
				m.Lock()
				delete(Data[channel].Clients, connId)
				m.Unlock()
				conn.Close()
				break loop
			}
		}
	}
}
func addClient(id, channel int, addr string) {
	m.Lock()
	if copy, ok := Data[channel]; ok {
		copy.Clients[id] = addr
		Data[channel] = copy
	}
	m.Unlock()
}
