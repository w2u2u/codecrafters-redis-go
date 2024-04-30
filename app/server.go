package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/server"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	keyValueDb := database.NewKeyValue()
	redisServer := server.NewServer("6379", &keyValueDb)

	if err := redisServer.Master(); err != nil {
		fmt.Println("Unable to run the server:", err.Error())
		os.Exit(1)
	}
}
