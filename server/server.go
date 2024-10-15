package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

// represents the chat server
type Server struct {
	listenAdd  string              // Address and port to listen on
	ln         net.Listener        // TCP Listener
	clients    map[net.Conn]string // Connected clients
	muxLock    sync.Mutex          // Mutex ensures that only one goroutine can access a critical section of code at any given time.
	message    []string            // Previous messages
	quit       chan struct{}       // Channel to signal termination(server to shut down)
	maxClients int                 // Maximum number of clients
}

// Initializes a new server with the given listen address
func NewServer(listenAdd string) *Server {
	return &Server{
		listenAdd:  listenAdd,
		clients:    make(map[net.Conn]string),
		message:    []string{},
		quit:       make(chan struct{}),
		maxClients: 10,
	}
}

// Start the server and listen for incoming connections
func (s *Server) Start() error {
	// Create a TCP listener on the specified address
	ln, err := net.Listen("tcp", s.listenAdd)
	if err != nil {
		return err // Return error if unable to listen
	}
	defer ln.Close()

	fmt.Printf("Listening to connection on port: %s\n", s.listenAdd)
	s.ln = ln // Store the listener to the server instance

	go s.acceptLoop() // Start accepting connections

	<-s.quit // Wait for server to shut down
	return nil
}

// Handles incoming connections
func (s *Server) acceptLoop() {
	for {
		// Accept a new connection
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("Accept error: ", err) // Log any error during accept
			continue
		}

		// Lock to safely handle number of clients
		s.muxLock.Lock()
		if len(s.clients) >= s.maxClients {
			conn.Write([]byte("Server is full.\n"))
			s.muxLock.Unlock()
			conn.Close()
			continue
		}
		s.muxLock.Unlock()

		go s.handleClient(conn) // Handle new clients
	}
}

// Manges the communication between clients
func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close() // Close the connection when done

	// Get the name of the client
	name, err := s.getName(conn)
	if err != nil {
		return
	}

	// Lock to add the client to the list of connected clients
	s.muxLock.Lock()
	s.clients[conn] = name // Add the client to the list
	s.muxLock.Unlock()

	// Broadcast a message to all connected clients
	s.broadcast(fmt.Sprintf("%s has joined the chat...\n", name), conn)
	s.sendPreviousMessages(conn)
	s.promptAllClients()

	// Create a scanner to read messages from the client
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue // Skip the input if empty
		}

		// Format the message, stores it and broadcasts it
		formattedMessage := formatMessage(name, text)
		s.message = append(s.message, formattedMessage)
		s.broadcast(formattedMessage, conn)

		s.promptAllClients() // Update prompt
	}

	// Remove the client from the list of connected clients
	s.muxLock.Lock()
	delete(s.clients, conn) // Remove the client
	s.muxLock.Unlock()

	// Notify others that the client has left
	s.broadcast(fmt.Sprintf("%s has left the chat...\n", name), conn)
	s.promptAllClients()
}

// Get the name of the client
func (s *Server) getName(conn net.Conn) (string, error) {
	conn.Write([]byte(welcomeMessage(s.getClientsName())))
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		conn.Write([]byte("Invalid name. Disconnection....\n"))
		return "", err
	}
	return strings.TrimSpace(name), nil
}

// Broadcast a message to all connected clients except the sender
func (s *Server) broadcast(msg string, sender net.Conn) {
	s.muxLock.Lock()
	defer s.muxLock.Unlock()

	for conn := range s.clients {
		if conn != sender { // Avoids broadcasting to the sender
			if _, err := conn.Write([]byte("\n" + msg)); err != nil {
				log.Println("Error writing to connection: ", err)
			}
		}
	}
}

// Send a prompt to all connected clients
func (s *Server) promptAllClients() {
	s.muxLock.Lock()
	defer s.muxLock.Unlock()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	for conn, name := range s.clients {
		prompt := fmt.Sprintf("[%s][%s]:", timestamp, name)
		if _, err := conn.Write([]byte(prompt)); err != nil {
			log.Println("Error writing to connection: ", err)
		}
	}
}

// Sends all the previous messages to the new client
func (s *Server) sendPreviousMessages(conn net.Conn) {
	for _, msg := range s.message {
		conn.Write([]byte(msg))
	}
}

func formatMessage(name, text string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:%s\n", timestamp, name, text)
}

func welcomeMessage(existingUsers []string) string {
	userList := "Current users: "
	if len(existingUsers) > 0 {
		userList += strings.Join(existingUsers, ", ") + "\n"
	} else {
		userList += "None\n"
	}
	return userList + "Welcome to TCP-Chat!\n" +
		"         _nnnn_\n" +
		"        dGGGGMMb\n" +
		"       @p~qp~~qMb\n" +
		"       M|@||@) M|\n" +
		"       @,----.JM|\n" +
		"      JS^\\__/  qKL\n" +
		"     dZP        qKRb\n" +
		"    dZP          qKKb\n" +
		"   fZP            SMMb\n" +
		"   HZM            MMMM\n" +
		"   FqM            MMMM\n" +
		" __| \".        |\\dS\"qML\n" +
		"|    .       | ' \\Zq\n" +
		"_)      \\.___.,|     .'\n" +
		"\\____   )MMMMMP|   .'\n" +
		"     -'       --'\n" +
		"[ENTER YOUR NAME]:"
}

func (s *Server) getClientsName() []string {
	s.muxLock.Lock()
	defer s.muxLock.Unlock()
	names := make([]string, 0, len(s.clients))
	for _, name := range s.clients {
		names = append(names, name)
	}
	return names
}
