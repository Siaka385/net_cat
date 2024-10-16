package servers

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	mu.Lock()
	clientCount++

	conn.Write([]byte(WelcomeMessage()))

	newClient := Client{
		clientconn:  conn,
		clientsName: conn.RemoteAddr().String(),
	}
	clients[conn] = newClient

	mu.Unlock()

	scanner := bufio.NewScanner(newClient.clientconn)

	if scanner.Scan() {
		newClient.clientsName = strings.TrimSpace(scanner.Text())
		mu.Lock()
		clients[conn] = newClient
		mu.Unlock()
	}

	fmt.Printf("Client connected: %s (Total clients: %d)\n", newClient.clientsName, clientCount)

	broadcastMessage(conn, "\n"+newClient.clientsName+" has joined the chat...\n")
	if clientCount > 1 && LoadPreviusMessages {
		conn.Write([]byte(LoadPreviousChats() + "\n"))
		LoadPreviusMessages = false
	}

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
	defer func() {
		DisconnectChan <- conn
	}()
}

func HandleDisconnection(conn net.Conn) {
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

func HandleServerFull(conn net.Conn) {
	mu.Lock()
	defer mu.Unlock()

	conn.Write([]byte("The server is currently full\n"))
	conn.Close()
}
