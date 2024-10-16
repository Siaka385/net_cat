package servers

import (
	"fmt"
	"net"
)

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
