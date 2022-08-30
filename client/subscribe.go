package client

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

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

		// It is possible that no file is being transmitted
		if noFile(conn, string(data[:n])) {
			continue
		}

		finfo := strings.Split(string(data[:n]), "|")
		fsize, _ := strconv.ParseInt(finfo[1], 10, 64)

		// It could also happen that the file being transmitted is already in the client's system
		// we ask the user if he wants it anyway.
		if !wantsFile(finfo[0]) {
			if _, err := conn.Write([]byte("No")); err != nil {
				fmt.Println("Error sending No", err)
				os.Exit(1)
			}
			continue
		}
		fmt.Printf("Server is transmitting: %v\n", finfo[0])

		if _, err := conn.Write([]byte("OK")); err != nil {
			fmt.Println("Error sending OK", err)
			os.Exit(1)
		}

		f, err := os.Create(finfo[0])
		if err != nil {
			fmt.Printf("Error creating %s %v\n", finfo[0], err)
			os.Exit(1)
		}

		written, err := io.CopyN(f, conn, fsize)
		if err != nil {
			fmt.Printf("Error copying %d/%d data to file %v", written, fsize, err)
			os.Exit(1)
		}

		fmt.Printf("Received %d bytes from server, copied to %v\n", written, finfo[0])
		f.Close()
		conn.Write([]byte("ok"))
	}
}

// noFile, is called in case no file is being transmitted on the specified channel.
// The program asks the user if he wants to disconnect, if that's the case then a
// disconnect signal is sent to the server and the client exits.
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

// wantsFile makes sure the user wants the file being transmitted even if said file already exists.
func wantsFile(filename string) bool {
	_, err := os.OpenFile(filename, os.O_RDONLY, 4440)
	return errors.Is(err, os.ErrNotExist)
}
