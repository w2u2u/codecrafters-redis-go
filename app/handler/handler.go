package handler

import (
	"bufio"
	"encoding/hex"
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
					fmt.Sprintf("master_replid:%s", handler.cfg.MasterReplid),
					"master_repl_offset:0",
				}
				handler.writer.WriteString(command.NewBulkString(strings.Join(info, "\n")))
			}

		case command.Replconf:
			if len(redisCmd.Args) < 3 {
				fmt.Println("Replconf command requires arguments")
				break cmdLoop
			}

			handler.writer.WriteString(command.NewSimpleString("OK"))

		case command.Psync:
			if len(redisCmd.Args) < 3 {
				fmt.Println("Psync command requires arguments")
				break cmdLoop
			}

			resp := fmt.Sprintf("FULLRESYNC %s 0", handler.cfg.MasterReplid)
			handler.writer.WriteString(command.NewSimpleString(resp))

			data, err := hex.DecodeString("524544495330303131fa0972656469732d76657205372e322e30fa0a72656469732d62697473c040fa056374696d65c26d08bc65fa08757365642d6d656dc2b0c41000fa08616f662d62617365c000fff06e3bfec0ff5aa2")
			if err != nil {
				fmt.Println("Unable to decode hex to bytes:", err)
				break cmdLoop
			}
			resp = fmt.Sprintf("$%d\r\n%s", len(data), data)
			handler.writer.WriteString(resp)
		}

		handler.writer.Flush()
	}
}

func (handler *Handler) Handshake() error {
	handler.writer.WriteString(command.NewArrays([]string{"PING"}))
	handler.writer.Flush()

	_, err := handler.reader.ReadString('\n')
	if err != nil {
		return err
	}

	handler.writer.WriteString(command.NewArrays([]string{
		"REPLCONF",
		"listening-port",
		handler.cfg.Port,
	}))
	handler.writer.Flush()

	_, err = handler.reader.ReadString('\n')
	if err != nil {
		return err
	}

	handler.writer.WriteString(command.NewArrays([]string{
		"REPLCONF",
		"capa",
		"psync2",
	}))
	handler.writer.Flush()

	_, err = handler.reader.ReadString('\n')
	if err != nil {
		return err
	}

	handler.writer.WriteString(command.NewArrays([]string{
		"PSYNC",
		"?",
		"-1",
	}))
	handler.writer.Flush()

	_, err = handler.reader.ReadString('\n')
	if err != nil {
		return err
	}

	return nil
}
