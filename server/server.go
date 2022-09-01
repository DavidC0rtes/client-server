package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

const (
	CONN_HOST = "127.0.0.1"
	CONN_PORT = "3000"
	CONN_TYPE = "tcp"
)

// Info Metadata
type Info struct {
	CurrFile string
	Filesize int64
	Total    int64
	Clients  map[int]string
	MaxSize  int64
}

var chans []chan []byte
var quit = make(chan int)
var done = make(chan bool, 1)
var Data = make(map[int]Info)

var m sync.Mutex

var MAX_SIZE int64

// Run, starts the server with the desired number of channels and desired max filesize.
func Run(numChannels int, maxFilesize int64) {
	fmt.Printf("Server starting with %d channels and max file size of %d(B)\n", numChannels, maxFilesize)

	// Create and initialize every channel and the Data struct.
	for i := 0; i < numChannels; i++ {
		chans = append(chans, make(chan []byte))

		Data[i] = Info{
			"",
			0,
			0,
			make(map[int]string),
			maxFilesize,
		}
	}
	MAX_SIZE = maxFilesize

	listen, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer listen.Close()
	fmt.Printf("Waiting for incoming requests on %s:%s\n", CONN_HOST, CONN_PORT)

	go startAPI()
	// forever...
	for i := 0; ; i++ {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(2)
		}
		// i serves as an id for the connection, this will be useful
		// to know which client disconnects.
		go handleIncomingRequest(conn, i)
	}
}

func handleIncomingRequest(conn net.Conn, id int) {
	// store incoming data
	buffer := make([]byte, 4096)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error reading protocol message")
		return
	}

	processRequest(string(buffer[:n]), conn, id)
	conn.Close()
}

// Process any given request sent to the server.
// A valid request has the form:
// Sending files: -> <content-size> <file> <channel>
// Subscribing to channel: listen <channel>
func processRequest(body string, conn net.Conn, id int) {

	fmt.Printf(">>>>>> %v\n", body)

	content, err := SplitRequest(body)
	if err != nil {
		fmt.Println(err)
		return
	}

	channel, _ := strconv.Atoi(content[len(content)-1])
	clientAddr := conn.RemoteAddr().String()

	switch {
	case content[0] == "->": // If we are receiving from a client...

		addClient(id, channel, clientAddr)

		// Every access to Data needs to be inside a mutual exclusion block
		// bc it is a global variable and is not thread-safe
		go func() {
			<-done
			//m.Lock()
			//fmt.Println("goodbye")
			delete(Data[channel].Clients, id)
			//m.Unlock()
			return
		}()

		receiveFile(content[1], content[2], channel, id, conn)

	case content[0] == "listen": // If a client wants to listen...
		addClient(id, channel, clientAddr)
		sendtoClient(channel, id, conn)

	}
}

func addClient(id, channel int, addr string) {
	m.Lock()
	if copy, ok := Data[channel]; ok {
		copy.Clients[id] = addr
		Data[channel] = copy
	}
	m.Unlock()
}

// SplitRequest returns a []string of every part of the request.
// But first makes sure said request comply with the protocol.
func SplitRequest(req string) ([]string, error) {
	matchListen, err := regexp.MatchString(`^(listen)\s+\d+$`, req)
	if err != nil {
		return nil, err
	}

	regexSend, err := regexp.Compile(`^(->)\s+\d+\s+([a-zA-Z0-9-_])+(\.[a-zA-Z0-9]+)?\s+\d+$`)
	if err != nil {
		return nil, err
	}

	matchSend := regexSend.MatchString(req)

	if !matchListen && !matchSend {
		return nil, errors.New("Malformed request")
	}
	return strings.Fields(req), nil
}
