package servers

import (
	"net"
	"testing"
	"time"
)

// Checking if a connection is able to connect to the servrt
func TestConnect(t *testing.T) {
	port := ":8080"

	// start the server
	go Connection(port)

	// give time for the server to start
	// Sleep pauses the current goroutine for at least the duration d. A negative or zero duration causes Sleep to return immediately.
	time.Sleep(1 * time.Second)

	// create a client to connect to the server
	conn, err := net.Dial("tcp", port)
	if err != nil {
		t.Error("failed to set up the client", err)
	}
	// close the client connection before the function returns
	defer conn.Close()

	// check if client is able to send a message to the server
	_, err = conn.Write([]byte("Hello server"))
	if err != nil {
		t.Error("failed to send messsage to the server", err)
	}
}

// // Initialize shared variables for testing
// func setupTestEnv() {
// 	mu = sync.Mutex{}
// 	clients = make(map[net.Conn]Client)
// 	clientCount = 0
// 	UserMessages = []string{}
// 	firstmessage = true
// 	LoadPreviusMessages = true
// }

// func TestHandleConnection(t *testing.T) {
// 	// Setup test environment (reset shared state)
// 	setupTestEnv()

// 	// Create a pipe to simulate a TCP connection (net.Pipe provides a pair of connected sockets)
// 	serverConn, clientConn := net.Pipe()
// 	defer serverConn.Close()
// 	defer clientConn.Close()

// 	// Run the handleConnection function in a goroutine to handle the connection asynchronously
// 	go handleConnection(serverConn)

// 	// Simulate client behavior: send the client's name
// 	clientWriter := bufio.NewWriter(clientConn)
// 	// clientReader := bufio.NewReader(clientConn)

// 	// Simulate sending client's name
// 	clientWriter.WriteString("TestClient\n")
// 	clientWriter.Flush()

// 	// Give the server a moment to process the connection
// 	time.Sleep(100 * time.Millisecond)

// 	// Check that the client's name has been registered correctly
// 	mu.Lock()
// 	if clientCount != 1 {
// 		t.Errorf("Expected 1 client, got %d", clientCount)
// 	}

// 	// Verify client registration
// 	client, ok := clients[serverConn]
// 	if !ok {
// 		t.Fatal("Expected client to be registered, but it's not found")
// 	}
// 	if client.clientsName != "TestClient" {
// 		t.Errorf("Expected client name to be 'TestClient', got '%s'", client.clientsName)
// 	}
// 	mu.Unlock()

// 	// Simulate the client sending a message
// 	clientWriter.WriteString("Hello, World!\n")
// 	clientWriter.Flush()

// 	// Give the server a moment to process the message
// 	time.Sleep(100 * time.Millisecond)

// 	// Simulate client disconnection
// 	clientConn.Close()
// 	time.Sleep(100 * time.Millisecond)

// 	// Check that the client count has been decremented
// 	mu.Lock()
// 	if clientCount != 0 {
// 		t.Errorf("Expected 0 clients after disconnect, got %d", clientCount)
// 	}
// 	mu.Unlock()
// }
