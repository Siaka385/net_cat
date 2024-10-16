package servers

import (
	"fmt"
	"net"
	"sync"
)

type Client struct {
	clientconn  net.Conn
	clientsName string
}

var (
	clients             = make(map[net.Conn]Client) // map to hold clients and
	mu                  sync.Mutex                  // mutex to synchronize access to the clients map
	clientCount         int                         // variable to count connected clients
	DisconnectChan      = make(chan net.Conn)
	UserMessages        []string
	firstmessage        bool // to avoid sending empty message on first instance of creations
	LoadPreviusMessages bool
)

func Connection(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error")
	}

	defer listener.Close()

	fmt.Println("open connection on port", port)
	for {

		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error")
			continue
		}
		firstmessage = true
		LoadPreviusMessages = true
		go handleConnection(conn)
	}
}
