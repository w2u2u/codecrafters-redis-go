package handler

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/database"
)

type Handler struct {
	db     database.IDatabase
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewHandler(conn net.Conn, db database.IDatabase) Handler {
	return Handler{
		db:     db,
		conn:   conn,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (handler *Handler) Handle() {
	defer handler.conn.Close()

cmdLoop:
	for {
		redisCommand, err := command.NewRedisCommand(handler.reader)
		if err != nil {
			fmt.Println("Unable to parse redis command:", err)
			break
		}

		fmt.Printf("Redis command: %v\n", redisCommand)

		switch strings.ToLower(redisCommand.Args[0]) {

		case command.Ping:
			handler.writer.WriteString("+PONG\r\n")

		case command.Echo:
			if len(redisCommand.Args) < 2 {
				fmt.Println("Echo command requires an argument")
				break cmdLoop
			}

			echoArg := redisCommand.Args[1]

			handler.writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(echoArg), echoArg))

		case command.Get:
			if len(redisCommand.Args) < 2 {
				fmt.Println("Get command requires an argument")
				break cmdLoop
			}

			keyArg := redisCommand.Args[1]

			value, err := handler.db.Get(keyArg)
			if err != nil {
				handler.writer.WriteString("$-1\r\n")
				break
			}

			handler.writer.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))

		case command.Set:
			if len(redisCommand.Args) < 3 {
				fmt.Println("Set command requires arguments")
				break cmdLoop
			}

			keyArg, valueArg := redisCommand.Args[1], redisCommand.Args[2]
			expArg := "0"

			if len(redisCommand.Args) > 4 {
				unit := strings.ToLower(redisCommand.Args[3])
				expiry := redisCommand.Args[4]
				switch unit {
				case "ex":
					expArg = expiry + "s"
				case "px":
					expArg = expiry + "ms"
				}
			}

			handler.db.Set(keyArg, valueArg, expArg)
			handler.writer.WriteString("+OK\r\n")
		}

		handler.writer.Flush()
	}
}
