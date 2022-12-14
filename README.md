# File share 
CLI app to send and receive files between two or more clients and a server.
Clients can subscribe(listen) to specific channels and
servers coordinate them (receive files and redirect them to their destination).

## Custom protocol over TCP
A simple protocol to enable the communication between client and server.
* Sending files: `-> <content-size> <file> <channel>` 
* Subscribing to channel: `listen <channel>`

![diagram for the protocol](https://i.ibb.co/dffVj5J/arquitecfileshare-protocol-drawio.png)

### General working of the server
Click the image

<a href="https://ibb.co/6W32kcp"><img src="https://i.ibb.co/6W32kcp/flow-server.png" alt="flow-server" border="0"></a>

## Usage
```
Usage:
  client-server [command]

Available Commands:
  client      Initialize a client
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  server      Initialize the server

Flags:
  -h, --help     help for client-server
  -t, --toggle   Help message for toggle
```
### Starting the server
To start the server with default options run `./client-server server`. 
More options/flags:
```
Flags:
  -c, --channels int   Tells the number of channels to create. (default 3)
  -h, --help           help for server
  -m, --max int        Maximum supported filesize (B). (default 4096)
```
### Sending files 
To send a file to the server on channel e.g (1) run: 

`./client-server client send --channel 1 ../foo.txt`

To see additional flags and information run 

`./client-server client send --help`
### Listening(subscribing) to a channel
run: 
`./client-server client receive --channel 1` 

as usual help for the command is displayed by doing

`./client-server client --help`

## Frontend
The frontend serves as a reporting page where general information about the server is displayed:
* Current file and size being transmitted on each channel.
* Current clients listening over each channel.
* Amount of bytes transmitted since the last update on each channel.
* Line graph showing historic data of total concurrent clients.
* Filetypes transmitted.

The frontend running on port `5173` gets all this data from an API running on port `8080` by making a GET request to the `/info` route. Said API is started alongside the server but they are different applications.


## Built with
* net package.
* [gin](https://github.com/gin-gonic/gin) (API).

<img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" height="100px">

* [cobra](https://cobra.dev/) (CLI arguments).

<img src="https://cobra.dev/home/logo.png" height="100px">

* Vue.js (Front-end).

![Vue logo](https://upload.wikimedia.org/wikipedia/commons/thumb/9/95/Vue.js_Logo_2.svg/100px-Vue.js_Logo_2.svg.png)

* Charts.js

<img src="https://www.chartjs.org/img/chartjs-logo.svg" height="100px">

* Bootstrap 5 (styling and css classes).

![bootstrap logo](https://upload.wikimedia.org/wikipedia/commons/thumb/b/b2/Bootstrap_logo.svg/100px-Bootstrap_logo.svg.png)
