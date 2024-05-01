package handler

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/database"
)

type Handler struct {
	db     database.IDatabase
	conn   net.Conn
	cfg    config.Config
	reader *bufio.Reader
	writer *bufio.Writer
}

func NewHandler(conn net.Conn, db database.IDatabase, cfg config.Config) Handler {
	return Handler{
		db:     db,
		conn:   conn,
		cfg:    cfg,
		reader: bufio.NewReader(conn),
		writer: bufio.NewWriter(conn),
	}
}

func (handler *Handler) Handle() {
	defer handler.conn.Close()

cmdLoop:
	for {
		redisCmd, err := command.NewRedisCommand(handler.reader)
		if err != nil {
			fmt.Println("Unable to parse redis command:", err)
			break
		}

		fmt.Printf("Redis command: %v\n", redisCmd)

		switch strings.ToLower(redisCmd.Args[0]) {

		case command.Ping:
			handler.writer.WriteString(command.NewSimpleString("PONG"))

		case command.Echo:
			if len(redisCmd.Args) < 2 {
				fmt.Println("Echo command requires an argument")
				break cmdLoop
			}

			echoArg := redisCmd.Args[1]

			handler.writer.WriteString(command.NewBulkString(echoArg))

		case command.Get:
			if len(redisCmd.Args) < 2 {
				fmt.Println("Get command requires an argument")
				break cmdLoop
			}

			keyArg := redisCmd.Args[1]

			value, err := handler.db.Get(keyArg)
			if err != nil {
				handler.writer.WriteString(command.NewNulls())
				break
			}

			handler.writer.WriteString(command.NewBulkString(value))

		case command.Set:
			if len(redisCmd.Args) < 3 {
				fmt.Println("Set command requires arguments")
				break cmdLoop
			}

			keyArg, valueArg := redisCmd.Args[1], redisCmd.Args[2]
			expArg := "0"

			if len(redisCmd.Args) > 4 {
				unit := strings.ToLower(redisCmd.Args[3])
				expiry := redisCmd.Args[4]
				switch unit {
				case "ex":
					expArg = expiry + "s"
				case "px":
					expArg = expiry + "ms"
				}
			}

			handler.db.Set(keyArg, valueArg, expArg)
			handler.writer.WriteString(command.NewSimpleString("OK"))

		case command.Info:
			if len(redisCmd.Args) < 2 {
				fmt.Println("Info command requires arguments")
				break cmdLoop
			}

			sectionArg := redisCmd.Args[1]

			if sectionArg == "replication" {
				info := []string{
					fmt.Sprintf("role:%s", handler.cfg.Role),
					"connected_slaves:0",
					"master_replid:8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
				}
				handler.writer.WriteString(command.NewBulkString(strings.Join(info, "\n")))
			}
		}

		handler.writer.Flush()
	}
}
