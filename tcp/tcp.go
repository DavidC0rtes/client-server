package tcp

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "3000"
	CONN_TYPE = "tcp"
)

// Info Metadata
type Info struct {
	channel  chan []byte
	currFile string
	filesize int64
}

var channels = map[int]Info{
	0: {
		make(chan []byte),
		"",
		0,
	},
	1: {
		make(chan []byte),
		"",
		0,
	},
}

func Run() {
	fmt.Println("Server running...")

	listen, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer listen.Close()
	fmt.Printf("Waiting for incoming requests on %s:%s\n", CONN_HOST, CONN_PORT)

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
		os.Exit(1)
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

	switch {
	case content[0] == "->":
		channel, _ := strconv.Atoi(content[3])
		receiveFile(content[1], content[2], channel, conn)
	case content[0] == "listen":
		channel, _ := strconv.Atoi(content[1])
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

	if copy, ok := channels[channel]; ok {
		copy.currFile = filename
		copy.filesize = fileSize
		channels[channel] = copy
	}

	fmt.Printf("Emitting data over channel %d\n", channel)
	inputBuffer := make([]byte, fileSize)
	nBytes, err := conn.Read(inputBuffer)
	if err != nil {
		fmt.Println("Error reading file", err.Error())
	}
	// Emit forever on channel.
	fmt.Printf("%v %d\n", channels[channel].currFile, nBytes)
	// Hacer esto cuando se este enviando al cliente, cuando esté conectado.
	msg := fmt.Sprintf("%v %d", channels[channel].currFile, nBytes)
	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	for {
		channels[channel].channel <- inputBuffer
	}
}

/*
Sends a file to the clients listening on the specified channel.
*/
func sendtoClient(channel int, conn net.Conn) {
	givenChannel, ok := channels[channel]
	// Clients in receive mode aren't allowed to create channels.
	if !ok {
		fmt.Printf("Channel %d does not exist.\n")
		return
	}

	fmt.Printf("Subscribing to %d!\n", channel)

	// Receives and broadcast file contents to the channel.
	for {
		select {
		case data := <-givenChannel.channel:
			if _, err := conn.Write(data); err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}

}
