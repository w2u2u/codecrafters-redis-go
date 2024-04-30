package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/server"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	redisServer := server.NewServer("6379")

	if err := redisServer.Master(); err != nil {
		fmt.Println("Unable to run the server:", err.Error())
		os.Exit(1)
	}
}
