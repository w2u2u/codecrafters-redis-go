package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	cfg := config.Parse()
	keyValueDb := database.NewKeyValue()
	redisServer := server.NewServer(cfg, &keyValueDb)

	if err := redisServer.TryHandshake(); err != nil {
		fmt.Println("Unable to handshake the master:", err.Error())
		os.Exit(1)
	}

	l, err := redisServer.Listen()
	if err != nil {
		fmt.Println("Unable to run the server:", err.Error())
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			continue
		}

		go redisServer.HandleConnection(conn)
	}
}
