package server

import (
	"bufio"
	"fmt"
	"net"
	"testing"
	"time"
)

func TestServer(t *testing.T) {
	listenAdd := "localhost:2525"
	server := NewServer(listenAdd)

	go func() {
		if err := server.Start(); err != nil {
			t.Fatalf("Server start error: %v", err)
		}
	}()

	// Allow the server a moment to start
	time.Sleep(1 * time.Second)

	// Create clients
	clientCount := 3
	for i := 0; i <= clientCount; i++ {
		go func(clientID int) {
			conn, err := net.Dial("tcp", listenAdd)
			if err != nil {
				t.Fatalf("Failed to connect: %v", err)
			}
			defer conn.Close()

			// Send name
			fmt.Fprintf(conn, "Client %d\n", clientID)

			// Read messages from the server
			scanner := bufio.NewScanner(conn)
			for scanner.Scan() {
				message := scanner.Text()
				fmt.Printf("Client %d received: %s\n", clientID, message)
			}
		}(i)
	}

	// Wait for the server to finish
	time.Sleep(5 * time.Second)
	
	// Close the server
	close(server.quit)
}