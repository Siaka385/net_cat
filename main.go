package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"net_cafe.com/m/server"
)

func main() {
	port := "2525"

	if len(os.Args) > 0 && len(os.Args) < 3 {
		if len(os.Args) == 2 {
			_, err := strconv.ParseInt(os.Args[1], 10, 64)
			if err != nil {
				fmt.Println("[USAGE]: ./TCPChat $port")
				return
			}
			port = os.Args[1]
		}
	} else {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}

	Server := server.NewServer(":" + port)
	if err := Server.Start(); err != nil {
		log.Fatal("Error starting server: ", err)
	}
}
