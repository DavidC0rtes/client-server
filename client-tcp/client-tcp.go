package client_tcp

import (
	"fmt"
	"io"
	"net"
	"os"
)

func connect() net.Conn {
	// Connect to server
	conn, err := net.Dial("tcp", "127.0.0.1:3000")

	if err != nil {
		fmt.Println("Error connecting to server", err)
		os.Exit(1)
	}
	return conn
}

func PrepareSend(filename string, channel int) {
	conn := connect()

	fstat, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Communicate with server
	message := fmt.Sprintf("-> %d %s %d", fstat.Size(), fstat.Name(), channel)
	_, err = conn.Write([]byte(message))
	// Response from server
	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading response", err.Error())
		os.Exit(1)
	}
	// We know the server received the filename and filesize correctly,
	// so it is time to send the actual file.
	if s := string(buf[:n]); s == "OK" {
		sendFile(fstat, conn, buf)
	}
	conn.Close()
}

func sendFile(fstat os.FileInfo, conn net.Conn, buf []byte) {
	source, err := os.Open(fstat.Name())
	if err != nil {
		fmt.Println("Error opening file", err.Error())
		os.Exit(1)
	}

	nBytes, err := io.CopyN(conn, source, fstat.Size())
	if err != nil {
		fmt.Printf("Sent %d of %d bytes.\n ", nBytes, fstat.Size())
		os.Exit(1)
	}
	source.Close()

	fmt.Printf("Sent %d bytes to the server.\n", nBytes)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading response", err.Error())
		os.Exit(1)
	}
	fmt.Println(string(buf))
}

func Subscribe(channel int) {
	conn := connect()
	defer conn.Close()
	// Communicate with server
	message := fmt.Sprintf("listen %d", channel)
	_, err := conn.Write([]byte(message))
	if err != nil {
		fmt.Println("Error sending message", err.Error())
	}

	// Response from server
	ch := make(chan []byte)
	ech := make(chan error)
	go waitResponse(conn, ch, ech)

	//ticker := time.Tick(time.Second)
	for {
		select {
		// Received data from the connection
		case data := <-ch:
			fmt.Println(data)
		// Received an error  from the connection :(
		case err := <-ech:
			fmt.Println("Recieved error", err.Error())
			break
		}
	}
}

func waitResponse(conn net.Conn, ch chan []byte, che chan error) {
	for {
		// Read data
		data := make([]byte, 4096)
		_, err := conn.Read(data)
		if err != nil {
			che <- err
			return
		}
		ch <- data
	}
}
