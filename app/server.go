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

	if err := redisServer.Master(); err != nil {
		fmt.Println("Unable to run the server:", err.Error())
		os.Exit(1)
	}
}
