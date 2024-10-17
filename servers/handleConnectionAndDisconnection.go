package servers

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func handleConnection(conn net.Conn) {
	mu.Lock()     // Lock to safely modify shared resources
	clientCount++ // Increment the client count

	// Check if the server is full (more than 10 clients)
	if clientCount > 10 {
		conn.Write([]byte("server is currently full\n")) // Notify client that the server is full
		conn.Close()                                     // Close the connection
		return
	}

	conn.Write([]byte(WelcomeMessage())) // Send a welcome message to the client

	// Create a new client and store connection and remote address as default name
	newClient := Client{
		clientconn:  conn,
		clientsName: conn.RemoteAddr().String(), // Default client name is their remote address
	}
	clients[conn] = newClient // Add new client to the clients map

	mu.Unlock() // Unlock after modifying shared resources

	scanner := bufio.NewScanner(newClient.clientconn) // Create a scanner to read from client

	// If the client sends an initial message (their name)
	if scanner.Scan() {
		newClient.clientsName = strings.TrimSpace(scanner.Text()) // Set client's name
		mu.Lock()                                                 // Lock to update client info in map
		clients[conn] = newClient                                 // Update the client info in the map with the new name
		mu.Unlock()                                               // Unlock after updating
	}

	// Log the new client connection
	fmt.Printf("Client connected: %s (Total clients: %d)\n", newClient.clientsName, clientCount)

	// Broadcast the message that the new client has joined the chat
	broadcastMessage(conn, "\n"+newClient.clientsName+" has joined the chat...\n")

	// If there are other clients and the flag to load previous messages is true
	if clientCount > 1 && LoadPreviusMessages {
		// Check if there are previous chat messages
		if LoadPreviousChats() != "" {
			conn.Write([]byte(LoadPreviousChats() + "\n")) // Send previous messages to the new client
		}
		LoadPreviusMessages = false // Reset the flag after loading previous messages
	}

	// Handle incoming messages from this client
	handleMessage(newClient.clientconn, "", newClient.clientsName)

	// Infinite loop to continuously listen for messages from the client
	for {
		if scanner.Scan() {
			message := scanner.Text()                                           // Get the client's message
			handleMessage(newClient.clientconn, message, newClient.clientsName) // Process the client's message
		} else {
			break // Exit loop if there's no more input from the client
		}
	}

	// If there was an error reading from the client, log the error
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading from client %s: %v\n", newClient.clientsName, err)
	}

	defer func() {
		DisconnectChan <- conn // Notify that the client is disconnecting
	}()
}

func HandleDisconnection(conn net.Conn) {
	mu.Lock()         // Lock to safely modify shared resources
	defer mu.Unlock() // Unlock after modifications

	// Notify other clients that this client has left the chat
	go broadcastMessage(conn, "\n"+clients[conn].clientsName+" has left the chat...\n")

	// Check if the client exists in the map
	if client, ok := clients[conn]; ok {
		delete(clients, conn) // Remove the client from the clients map
		clientCount--         // Decrement the client count
		fmt.Printf("Client disconnected: %s (Total clients: %d)\n", client.clientsName, clientCount)
	}

	conn.Close() // Close the connection
}
