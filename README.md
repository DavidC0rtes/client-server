# File share 
App to send and receieve files between two or more clients and an arbitrer (server).
Clients suscribe(listen) to specific channels.
Servers coordinate them (receive files and redirect them).

## Custom protocol over TCP
A simple protocol to enable the communication between client and server.
* Sending files: `-> <data> <content-size> @<channel> <CLIENT-LOCAL-ADDRESS>` 
* Subscribing to channel: `listen <channel>`
* Client disconnecting from server: `disconnect <CLIENT-LOCAL-ADDRESS>`

The client's local address is used and saved for reporting statistics in the frontend.

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
To send a file to the server on channel e.g (1) run: `./client-server client send --channel 1 ../foo.txt`
To see additional flags and information run `./client-server client send --help`
### Listening(subscribing) to a channel
run: `./client-server client receive --channel 1`, as usual help for the command is displayed by appending `--help` after `receive`.

## Built with
* net package.
* [gin](https://github.com/gin-gonic/gin) (API).

<img src="https://raw.githubusercontent.com/gin-gonic/logo/master/color.png" height="100px">

* [cobra](https://cobra.dev/) (CLI arguments).

<img src="https://cobra.dev/home/logo.png" height="100px">

* Vue.js (Front-end).

![Vue logo](https://upload.wikimedia.org/wikipedia/commons/thumb/9/95/Vue.js_Logo_2.svg/100px-Vue.js_Logo_2.svg.png)

* Bootstrap 5 (styling and css classes).
	![bootstrap logo](https://upload.wikimedia.org/wikipedia/commons/thumb/b/b2/Bootstrap_logo.svg/100px-Bootstrap_logo.svg.png)
