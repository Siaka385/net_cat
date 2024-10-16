package main

import (
	"fmt"
	"os"

	"net_cat/servers"
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

	go servers.Connection(DefaultPort)

	for {
		select {
		case conn := <-servers.DisconnectChan:
			servers.HandleDisconnection(conn)
		}
	}
}
