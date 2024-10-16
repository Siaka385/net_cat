package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"net_cat/servers"
)

type Client struct {
	clientconn  net.Conn
	clientsName string
}

var (
	clients        = make(map[net.Conn]Client) // map to hold clients and
	mu             sync.Mutex                  // mutex to synchronize access to the clients map
	clientCount    int                         // variable to count connected clients
	disconnectChan = make(chan net.Conn)
	UserMessages   []string
	firstmessage   bool // to avoid sending empty message on first instance of creations
)

func main() {
	DefaultPort := ":2525"

	if len(os.Args[1:]) > 1 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
	if len((os.Args[1:])) == 1 {
		DefaultPort = ":" + os.Args[1:][0]
	}

	go connection(DefaultPort)

	for {
		select {
		case conn := <-disconnectChan:
			handleDisconnection(conn)
		}
	}
}

func connection(m string) {
	listener, err := net.Listen("tcp", m)
	if err != nil {
		fmt.Println("error")
	}

	defer listener.Close()

	fmt.Println("open connection on port", m)
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error")
		}
		firstmessage = true
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	mu.Lock()
	clientCount++

	conn.Write([]byte(servers.WelcomeMessage()))

	newClient := Client{
		clientconn:  conn,
		clientsName: conn.RemoteAddr().String(),
	}
	clients[conn] = newClient
	fmt.Printf("Client connected: %s (Total clients: %d)\n", newClient.clientsName, clientCount)
	mu.Unlock()

	defer func() {
		disconnectChan <- conn
	}()

	scanner := bufio.NewScanner(newClient.clientconn)

	if scanner.Scan() {
		newClient.clientsName = strings.TrimSpace(scanner.Text())
		mu.Lock()
		clients[conn] = newClient
		mu.Unlock()
	}
	broadcastMessage(conn, "\n"+newClient.clientsName+" has joined the chat...\n")

	handleMessage(newClient.clientconn, "", newClient.clientsName)

	for {
		if scanner.Scan() {
			message := scanner.Text()
			handleMessage(newClient.clientconn, message, newClient.clientsName)
		} else {
			break
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from client %s: %v\n", newClient.clientsName, err)
	}
}

func handleMessage(conn net.Conn, message, name string) {
	mu.Lock()

	defer mu.Unlock()
	client := clients[conn]
	conn.Write([]byte(MessageFormat(name) + ""))

	if firstmessage {
		firstmessage = false
		return
	}
	UserMessages = append(UserMessages, MessageFormat(client.clientsName)+message)

	go BroadcastUserMessage(conn, client.clientsName, message)
	fmt.Printf("Message from %s: %s\n", client.clientsName, message)
}

func BroadcastUserMessage(sender net.Conn, Sendername, Messages string) {
	mu.Lock()
	defer mu.Unlock()

	message := "\n" + MessageFormat(Sendername) + Messages
	for conn, client := range clients {
		if conn != sender {
			_, err := client.clientconn.Write([]byte(message + "\n" + MessageFormat(client.clientsName) + ""))
			if err != nil {
				fmt.Printf("Error broadcasting to %s: %v\n", client.clientsName, err)
			}
		}
	}
}

func handleDisconnection(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	go broadcastMessage(conn, "\n"+clients[conn].clientsName+" has left the chat...\n")

	if client, ok := clients[conn]; ok {

		delete(clients, conn)
		clientCount--
		fmt.Printf("Client disconnected: %s (Total clients: %d)\n", client.clientsName, clientCount)
	}

	conn.Close()
}

func MessageFormat(m string) string {
	times := time.Now()
	formartTime := times.Format("[2006-01-02 15:04:05]")

	return formartTime + "[" + m + "]" + ":"
}

func broadcastMessage(sender net.Conn, message string) {
	mu.Lock()
	defer mu.Unlock()
	for conn, client := range clients {
		if conn != sender {
			_, err := client.clientconn.Write([]byte(message + MessageFormat(client.clientsName) + ""))
			if err != nil {
				fmt.Printf("Error broadcasting to %s: %v\n", client.clientsName, err)
			}
		}
	}
}
