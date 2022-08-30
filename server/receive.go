package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
)

// Receives files from the connected client under the specified channel.
func receiveFile(size, filename string, channel, connId int, conn net.Conn) {
	// Convert size to int64
	fileSize, err := strconv.ParseInt(size, 10, 64)
	if err != nil {
		fmt.Println("Error reading filesize")
		return
	}

	if fileSize > MAX_SIZE {
		fmt.Printf("Error filesize (%d) exceeds maximum filesize allowed (%d)\n", fileSize, MAX_SIZE)
		return
	}

	if _, err = conn.Write([]byte("OK")); err != nil {
		fmt.Println("Couldn't send OK to client", err)
		return
	}

	m.Lock()
	if copy, ok := Data[channel]; ok {
		fmt.Println("H")

		dummyf, err := os.OpenFile(filename, os.O_RDONLY, 4440)
		if errors.Is(err, os.ErrNotExist) && Data[channel].CurrFile != "" {
			quit <- 1
		}
		if err == nil {
			dummyf.Close()
		}

		fmt.Println("D")
		copy.CurrFile = filename
		copy.Filesize = fileSize
		copy.Total = Data[channel].Total + fileSize

		Data[channel] = copy
	}
	m.Unlock()

	inputBuffer := make([]byte, fileSize)
	if _, err = conn.Read(inputBuffer); err != nil {
		fmt.Println("Error reading from input buffer", err)
		return
	}
	//done <- true
	fmt.Printf("Emitting data over channel %d\n", channel)
	for {
		select {
		case <-quit:
			fmt.Println("New file coming, stop")
			return
		default:
			fmt.Println("Sending")
			chans[channel] <- inputBuffer
		}
	}
}
