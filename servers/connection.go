package servers

import (
	"fmt"
	"net"
	"sync"
)

// Client represents a connected client with its connection and name
type Client struct {
	clientconn  net.Conn // The client's connection
	clientsName string   // The name or identifier of the client
}

var (
	clients             = make(map[net.Conn]Client) // Map to store clients with their connection as the key
	mu                  sync.Mutex                  // Mutex to protect shared resources and prevent race conditions across goroutines
	clientCount         int                         // Counter to track the number of connected clients
	DisconnectChan      = make(chan net.Conn)       // Channel to handle client disconnections
	UserMessages        []string                    // Slice to store messages from users
	firstmessage        bool                        // Flag to handle the first message sent (to avoid empty messages)
	LoadPreviusMessages bool                        // Flag to determine if previous messages should be loaded and sent to new clients
)

func Connection(port string) {
	// Listen for incoming TCP connections on the specified port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error") // Log error if there's an issue with starting the listener
	}

	defer listener.Close() // Ensure the listener is closed when done

	// Log that the server has opened the connection on the specified port
	fmt.Println("open connection on port", port)

	// Main loop to continuously accept new client connections
	for {
		conn, err := listener.Accept() // Accept a new connection
		if err != nil {
			fmt.Println("error") // Log error if connection acceptance fails
			continue             // Continue accepting new connections even if an error occurs
		}

		// Set flags for the first message and loading previous messages
		firstmessage = true
		LoadPreviusMessages = true

		// Handle the new client connection in a separate goroutine (concurrent handling)
		go handleConnection(conn)
	}
}
