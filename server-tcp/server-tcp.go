package server_tcp

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
	_, err = conn.Read(inputBuffer)
	if err != nil {
		fmt.Println("Error reading file", err.Error())
	}

	if err := os.WriteFile("out", inputBuffer, 0666); err != nil {
		fmt.Println("Couldn't create out file:", err)
		return
	}
	// Emit forever on channel.
	for {
		channels[channel].channel <- inputBuffer
	}
}

/*
Sends a file to the clients listening on the specified channel.
*/
func sendtoClient(channel int, conn net.Conn) {
	givenChannel, ok := channels[channel]
	if !ok {
		fmt.Printf("Channel %d does not exist.\n")
		return
	}

	fmt.Printf("Subscribing to %d\n", channel)

	// Receives and broadcast file contents to the channel.
	buf := make([]byte, 2)
	for {
		select {
		case data := <-givenChannel.channel:
			fmt.Printf("Sending %v\n", givenChannel.currFile)
			if _, err := conn.Write([]byte(givenChannel.currFile)); err != nil {
				fmt.Println("Couldn't send filename to client", err)
				return
			}

			if _, err := conn.Read(buf); err != nil {
				fmt.Println("Couldn't read OK from client", err)
				return
			}

			if string(buf) == "OK" {
				n, err := conn.Write(data)
				if err != nil {
					fmt.Println("Data not sent", err)
					return
				}
				fmt.Printf("Sent %dB to connection\n", n)
				/*hmm := strconv.Itoa(int(givenChannel.filesize))
				_, err := conn.Write([]byte(hmm))
				if _, err := conn.Read(buf); err != nil {
					fmt.Println("Couldn't read x from client", err)
					return
				}

				file, err := os.Open("out")
				if err != nil {
					fmt.Println("Couldn't open out", err)
					return
				}

				if n, err := io.CopyN(conn, file, givenChannel.filesize); err != nil {
					fmt.Printf("Copied %d bytes instead of %d %v", n, givenChannel.filesize, err)
					return
				}

				if err := file.Close(); err != nil {
					fmt.Println("Couldn't close file", err)
					return
				}
				fmt.Printf("File sent succesfully\n")*/
			}
		}
	}

}
