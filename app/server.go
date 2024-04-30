package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/command"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			os.Exit(1)
		}

		go handle(conn)
	}
}

func handle(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		redisCommand, err := command.NewRedisCommand(reader)
		if err != nil {
			fmt.Println("Unable to parse redis command:", err)
			break
		}

		fmt.Printf("Redis command: %v\n", redisCommand)

		switch strings.ToLower(redisCommand.Args[0]) {
		case command.Ping:
			writer.WriteString("+PONG\r\n")

		case command.Echo:
			if len(redisCommand.Args) < 2 {
				fmt.Println("Echo command requires an argument")
				break
			}
			echoArg := redisCommand.Args[1]
			writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(echoArg), echoArg))
		}

		writer.Flush()
	}
}
