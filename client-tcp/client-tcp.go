package client_tcp

import (
	"fmt"
	"io"
	"net"
	"os"
)

func Start(filename string, channel int) {
	fstat, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	// Connect to server
	conn, err := net.Dial("tcp", "127.0.0.1:3000")

	if err != nil {
		fmt.Println("Error connecting to server", err)
		os.Exit(1)
	}

	// Communicate with server
	message := fmt.Sprintf("-> %d %s %d\n", fstat.Size(), fstat.Name(), channel)
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
		source, err2 := os.Open(fstat.Name())
		if err2 != nil {
			fmt.Println("Error opening file", err2.Error())
			os.Exit(1)
		}
		nBytes, err1 := io.CopyN(conn, source, fstat.Size())
		if err1 != nil {
			fmt.Println("Error sending data to server", err1)
			os.Exit(1)
		}
		fmt.Printf("Sent %d bytes to the server.\n", nBytes)
		source.Close()
	}
	conn.Close()
}
