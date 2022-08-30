package client_tcp

import (
	"errors"
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

		if noFile(conn, string(data[:n])) {
			continue
		}

		name := string(data[:n])
		fmt.Printf("Received name %v %d\n", name, len(name))

		if !wantsFile(name) {
			if _, err := conn.Write([]byte("No")); err != nil {
				fmt.Println("Error sending No", err)
				os.Exit(1)
			}
			continue
		}

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

func disconnect(conn net.Conn) {
	msg := fmt.Sprintf("disconnect %s", conn.LocalAddr().String())
	fmt.Println("Sending disconnect!")
	if _, err := conn.Write([]byte(msg)); err != nil {
		fmt.Println("Error sending disconnect msg", err)
		os.Exit(1)
	}
}

func noFile(conn net.Conn, msg string) bool {
	if msg == "Nofile" {
		choice := "Y"
		fmt.Println("There are no files being shared on this channel. Disconnect? (Y/n)")
		fmt.Scanln(&choice)

		_, err := conn.Write([]byte(choice))
		if err != nil {
			fmt.Printf("Unable to send %s to server\n", choice)
		}
		fmt.Printf("Sent %s to server\n", choice)
		if choice != "n" {
			os.Exit(0)
		}
		return true
	}
	return false
}

func wantsFile(filename string) bool {
	_, err := os.OpenFile(filename, os.O_RDONLY, 4440)
	if !errors.Is(err, os.ErrNotExist) {
		fmt.Printf("%s is already in the system, do you want to receive it anyway? (y/N)\n", filename)
		choice := "N"
		fmt.Scanln(&choice)
		if choice != "y" {
			return false
		}
	}

	return true
}
