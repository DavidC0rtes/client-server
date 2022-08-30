package client_tcp

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type NewFile struct {
	name string
	size int64
}

var fileFromServer = NewFile{
	"",
	0,
}

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

	// We need to detect when ctrl c (SIGINT) is sent
	// to tell the server to disconnect the client
	// and update the Data struct.
	chanSignal := make(chan os.Signal)
	signal.Notify(chanSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-chanSignal
		disconnect(conn)
		os.Exit(0)
	}()

	return conn
}

// PrepareSend tells the server the name and size of the file
// it intends to share beforehand. This is important because
// the server needs to know how many bytes to expect and how to name
// the file.
func PrepareSend(filename string, channel int) {
	conn := connect()

	fstat, err := os.Stat(filename)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// Send special protocol message.
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
		fmt.Printf("Sent %dB of %dB.\n ", nBytes, size)
		os.Exit(1)
	}
	source.Close()

	fmt.Printf("Sent %dB to the server.\n", nBytes)
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
	data := make([]byte, 4096)

	for {
		// (1) Read name
		n, err := conn.Read(data)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		name := string(data[:n])
		fmt.Printf("Received name %v %d\n", name, len(name))

		if _, err := conn.Write([]byte("OK")); err != nil {
			fmt.Println("Error sending OK", err)
			os.Exit(1)
		}

		n1, err := conn.Read(data)
		if err != nil {
			os.Exit(1)
		}

		err = os.WriteFile(name, data[:n1], 0666)
		if err != nil {
			fmt.Printf("Couldn't write %v %v", name, err)
			os.Exit(1)
		}

		fmt.Printf("Received %d bytes from server copied to %v\n", n1, name)
		conn.Write([]byte("ok"))
	}
}

func waitResponse(conn net.Conn, ch chan []byte, che chan error) {
	data := make([]byte, 4096)
	for {
		fmt.Println("Hola")
		// Read data
		n, err := conn.Read(data)
		if err != nil {
			che <- err
			return
		}
		fmt.Printf("Receved name %v\n", string(data[:n]))
		conn.Write([]byte("OK"))

		n, err = conn.Read(data)
		if err != nil {
			che <- err
			return
		}
	}
}

func disconnect(conn net.Conn) {
	msg := fmt.Sprintf("disconnect %s", conn.LocalAddr().String())
	fmt.Println("Sending disconnect!")
	if _, err := conn.Write([]byte(msg)); err != nil {
		fmt.Println("Error sending disconnect msg", err)
		os.Exit(1)
	}
}
