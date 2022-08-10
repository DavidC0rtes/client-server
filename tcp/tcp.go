package tcp

import (
	"fmt"
	"io"
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
		fmt.Println("Error reading:", err.Error())
		os.Exit(1)
	}

	processRequest(string(buffer[:n]), conn)
	conn.Close()
}

/*
*
Sending files: -> <content-size> <file> @<channel>
Subscribing to channel: listen <channel>
*/
func processRequest(body string, conn net.Conn) {
	content := strings.Split(body, " ")

	switch {
	case content[0] == "->":
		receiveFile(content[1], content[2], conn)
	case content[0] == "listen":
		fmt.Println("Subscribing!")
	default:
		fmt.Printf("Malformed request: %s\n", body)
		os.Exit(1)
	}
}

func receiveFile(size, filename string, conn net.Conn) {
	// Convert size to int64
	fileSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		fmt.Println("Error reading file size")
		os.Exit(1)
	}

	if _, err = conn.Write([]byte("OK")); err != nil {
		fmt.Println("Error sending ok?", err.Error())
	}

	destFile, err := os.Create("out")
	if err != nil {
		fmt.Println("Error creating destination file: ", err.Error())
		os.Exit(1)
	}
	defer destFile.Close()

	// Time to receive file contents.
	n, err := io.Copy(destFile, conn)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	if n != fileSize {
		fmt.Printf("Filesize was %d, but %d bytes written\n", fileSize, n)
		return
	}
	fmt.Println("OI")
}
