package servers

import (
	"fmt"
	"net"
)

func handleMessage(conn net.Conn, message, name string) {
	mu.Lock()         // Lock to safely modify shared resources
	defer mu.Unlock() // Unlock once the function completes

	client := clients[conn] // Get the client associated with this connection

	// Send the formatted message header (name) to the client
	conn.Write([]byte(MessageFormat(name) + ""))

	// If it's the first message being handled, just return (do not broadcast)
	if firstmessage {
		firstmessage = false
		return
	}

	// If there's an actual message from the client (not an empty string)
	if message != "" {
		fmt.Println(MessageFormat(client.clientsName)) // Print the formatted message
		// Append the user's message to the global message store
		UserMessages = append(UserMessages, MessageFormat(client.clientsName)+message)
	}

	// Broadcast the user's message to all other clients
	go BroadcastUserMessage(conn, client.clientsName, message)

	// Log the message being handled
	fmt.Printf("Message from %s: %s\n", client.clientsName, message)
}

func BroadcastUserMessage(sender net.Conn, Sendername, Messages string) {
	mu.Lock()         // Lock to ensure thread-safe operations on shared resources
	defer mu.Unlock() // Unlock after broadcasting

	// Format the message with the sender's name and the message content
	message := "\n" + MessageFormat(Sendername) + Messages

	// Iterate over all connected clients
	for conn, client := range clients {
		// Skip broadcasting the message to the sender
		if conn != sender {
			// Attempt to write the message to the client's connection
			_, err := client.clientconn.Write([]byte(message + "\n" + MessageFormat(client.clientsName) + ""))
			if err != nil {
				// Log an error if there's an issue sending the message
				fmt.Printf("Error broadcasting to %s: %v\n", client.clientsName, err)
			}
		}
	}
}

func broadcastMessage(sender net.Conn, message string) {
	mu.Lock()         // Lock to ensure thread-safe operations on shared resources
	defer mu.Unlock() // Unlock after broadcasting

	// Iterate over all connected clients
	for conn, client := range clients {
		// Skip broadcasting to the sender
		if conn != sender {
			// Attempt to write the broadcast message to the client's connection
			_, err := client.clientconn.Write([]byte(message + MessageFormat(client.clientsName) + ""))
			if err != nil {
				// Log an error if there's an issue sending the broadcast message
				fmt.Printf("Error broadcasting to %s: %v\n", client.clientsName, err)
			}
		}
	}
}
