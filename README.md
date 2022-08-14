# Server app to send and receieve files between two or more clients.
Clients suscribe(listen) to specific channels.


## Custom protocol over TCP
A simple protocol to enable the communication between client and server.
* Sending files: `-> <data> <content-size> @<channel>` 
* Subscribing to channel: `listen <channel>`