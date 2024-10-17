# Net-cat

Net-cat is a Go-based application that recreates the functionality of the NetCat (nc) command in a server-client architecture. It allows multiple clients to connect via TCP, exchange messages in a group chat, and provides real-time updates when clients join or leave the chat.

## Features

1. TCP Server: Handles multiple clients (up to 10) simultaneously, allowing them to connect to the chat.
2. Client Naming: Each client is required to provide a name upon connection, which is used to identify messages.
3. Real-time Group Chat: Clients can send messages to the chat, which are broadcast to all connected clients, except empty messages.
4. Message Formatting: Messages include the timestamp and client name in the format: [YYYY-MM-DD HH:MM:SS][client_name]:[message]
5. Message History: When a client joins, they receive the chat history.
6. Join/Leave Notifications: All clients are notified when a new client joins or leaves the chat.
7. Default and Custom Ports: Server runs on port 8989 by default or a specified port.

## Usage
### Running the Server

- To run the server on the default port (8989):

```bash
$ go build

$./net_cat

open connection on port :8989
```
- To specify a different port:

```bash
$ go build

$./net_cat 2525

open connection on port :2525
```
- Connecting a Client

- Clients can connect using NetCat (nc) or another TCP client:

```bash

$ nc localhost 2525
```
- Upon connecting, the client is greeted with a welcome message and a Linux ASCII logo. The server then prompts for the client’s name. The chat begins once a name is provided.

Example session:

```bash

Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    `.       | `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     `-'       `--'

[ENTER YOUR NAME]: Alice
[2024-10-16 10:03:43][Alice]:Hello!
```
### Handling Errors

- If no port is specified or there’s a misuse of the command:

```bash
[USAGE]: ./net_cat $port
``` 

## Contributors
[Teddy Ogola Siaka](https://github.com/Siaka385)
