package server

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

// Sends a file to the clients listening on the specified channel.
func sendtoClient(channel, connId int, conn net.Conn) {

	if _, ok := Data[channel]; !ok {
		fmt.Printf("Channel %d does not exist.\n", channel)
		return
	}
	fmt.Printf("Subscribing to %d\n", channel)

loop:
	for {
		buf := make([]byte, 2)
		select {
		case data := <-chans[channel]:
			finfo := fmt.Sprintf("%s|%d", Data[channel].CurrFile, Data[channel].Filesize)
			fmt.Printf("Sending file info %v-%d\n", finfo, len(data))
			m.Lock()
			_, err := conn.Write([]byte(finfo))
			m.Unlock()
			if err != nil {
				fmt.Println("Couldn't send file info to client", err)
				break loop
			}

			fmt.Println("Waiting on OK")
			if _, err := conn.Read(buf); err != nil {
				fmt.Println("Couldn't read OK from client")
				break loop
			}

			if string(buf) == "OK" {
				//n, err = conn.Write(data)
				r := bytes.NewReader(data)
				//s := strings.NewReader(string(data))
				//fmt.Println(r.Read(data))
				nb, err := io.CopyN(conn, r, Data[channel].Filesize)
				if err != nil {
					fmt.Println("Couldn't send data to client", err)
					break loop
				}

				fmt.Printf("Sent %dB to connection\n", nb)

				if _, err := conn.Read(buf); err != nil {
					fmt.Println("Couldn't read response from client", err)
					break loop
				}
			}
		default:
			conn.Write([]byte("Nofile"))
			bf := make([]byte, 1)

			n, err := conn.Read(bf)
			if err != nil {
				fmt.Println("Couldn't read response from client", err)
				break
			}
			resp := string(bf[:n])
			if resp != "n" {
				fmt.Printf("Client %s whishes to disconnect.\n", conn.RemoteAddr().String())
				m.Lock()
				delete(Data[channel].Clients, connId)
				m.Unlock()
				conn.Close()
				break loop
			}

		}
	}
}
