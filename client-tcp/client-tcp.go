package client_tcp

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
)

/*
Starts a connection with the server on ::3000
*/
func connect() net.Conn {
	// Connect to server
	conn, err := net.Dial("tcp", "127.0.0.1:3000")

	if err != nil {
		fmt.Println("Error connecting to server", err)
		os.Exit(1)
	}
	return conn
}

// PrepareSend tells the server the name and size of the file
// it intends to share beforehand. This is important because
// the server needs to know how many bytes to expect and how to name
// the file. /*
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
		sendFile(filename, fstat.Size(), conn)
	}
	conn.Close()
}

// sendFile sends the actual file once we make sure the server
// got the name and size correctly.
func sendFile(filepath string, size int64, conn net.Conn) {
	source, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		os.Exit(1)
	}

	nBytes, err := io.CopyN(conn, source, size)
	if err != nil {
		fmt.Printf("Sent %d of %d bytes.\n ", nBytes, size)
		os.Exit(1)
	}
	source.Close()

	fmt.Printf("Sent %d bytes to the server.\n", nBytes)
}

/*
Subscribe clients can also subscribe/listen to a particular "channel"
this function sends that kind of request to the server and reacts to
responses.
*/
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

	for {
		select {
		// Received data from the connection
		case data := <-ch:
			// If there's a different file being transmitted or filesize changed...
			if _, err := os.Stat(string(data[:20])); err == nil {
				fmt.Println("File already in disk.")
				if _, err := conn.Write([]byte("NO")); err != nil {
					fmt.Println("Couldn't send NO to server", err)
					os.Exit(1)
				}

			} else if errors.Is(err, os.ErrNotExist) {

				recFile, err := os.Create(string(data[:20]))

				write, err := recFile.Write(data[21:])
				if err != nil {
					fmt.Println("Error writing file:", err)
					os.Exit(1)
				}
				fmt.Printf("Wrote %d bytes to %v.\n", write, recFile.Name())
			}

		// Received an error  from the connection :(
		case err := <-ech:
			fmt.Println("Received error", err.Error())
			os.Exit(1)
		}
	}
}

func waitResponse(conn net.Conn, ch chan []byte, che chan error) {
	data := make([]byte, 4096)
	for {
		// Read data
		n, err := conn.Read(data)
		if err != nil {
			che <- err
			return
		}
		ch <- data[:n]
	}
}
