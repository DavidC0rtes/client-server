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
	Clients  []string
}

var chans = []chan []byte{
	make(chan []byte),
	make(chan []byte),
}

var Data = map[int]Info{
	0: {
		"",
		0,
		0,
		make([]string, 100),
	},
	1: {
		"",
		0,
		0,
		make([]string, 100),
	},
}

var m sync.Mutex

func Run() {
	fmt.Println("Server running...")

	listen, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer listen.Close()
	fmt.Printf("Waiting for incoming requests on %s:%s\n", CONN_HOST, CONN_PORT)
	go startAPI()
	// forever...
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(2)
		}
		go handleIncomingRequest(conn)
	}
}

func handleIncomingRequest(conn net.Conn) {
	// store incoming data
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading request:", err.Error())
		return
	}

	processRequest(string(buffer[:n]), conn)
	conn.Close()
}

/*
Process any given request sent to the server.
A valid request has the form:

	Sending files: -> <content-size> <file> <channel>
	Subscribing to channel: listen <channel>
*/
func processRequest(body string, conn net.Conn) {
	content := strings.Split(body, " ")
	fmt.Printf(">>>>>> %v\n", body)

	channel, _ := strconv.Atoi(content[len(content)-2])
	clientAddr := content[len(content)-1]

	switch {
	case content[0] == "->":
		receiveFile(content[1], content[2], channel, conn)
	case content[0] == "listen":
		m.Lock()
		if copy, ok := Data[channel]; ok {
			copy.Clients = append(Data[channel].Clients, clientAddr)
			fmt.Printf("%v\n", copy.Clients)
			Data[channel] = copy
		}
		m.Unlock()

		sendtoClient(channel, conn)

	default:
		fmt.Printf("Malformed request: %s\n", body)
		os.Exit(1)
	}
}

/*
Receives files from the connected client under the specified channel.
*/
func receiveFile(size, filename string, channel int, conn net.Conn) {
	// Convert size to int64
	fileSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		fmt.Println("Error reading file size")
		os.Exit(1)
	}

	if _, err = conn.Write([]byte("OK")); err != nil {
		fmt.Println("Error sending OK signal to client", err.Error())
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
	if err != nil {
		fmt.Println("Error reading file", err.Error())
	}
	fmt.Printf("Emitting data over channel %d\n", channel)
	chans[channel] <- inputBuffer
}

/*
Sends a file to the clients listening on the specified channel.
*/
func sendtoClient(channel int, conn net.Conn) {

	if _, ok := Data[channel]; !ok {
		fmt.Printf("Channel %d does not exist.\n")
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
				fmt.Println("Couldn't read OK from client", err)
			}

			if string(buf) == "OK" {
				n, err := conn.Write(data)
				if err != nil {
					fmt.Println("Data not sent", err)
					return
				}
				fmt.Printf("Sent %dB to connection\n", n)
				if _, err := conn.Read(buf); err != nil {
					fmt.Println("Couldn't read ok from client", err)
					return
				}
			}
		}
	}
}

func GetData() *map[int]Info {
	return &Data
}
