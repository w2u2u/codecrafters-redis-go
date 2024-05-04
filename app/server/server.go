package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/command"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/database"
)

type Server struct {
	cfg    config.Config
	db     database.IDatabase
	slaves []net.Conn
}

func NewServer(cfg config.Config, db database.IDatabase) Server {
	return Server{
		cfg:    cfg,
		db:     db,
		slaves: make([]net.Conn, 0),
	}
}

func (s *Server) TryHandshake() error {
	if s.cfg.Role == "master" {
		return nil
	}

	conn, err := net.Dial("tcp", s.cfg.ReplicaOf)
	if err != nil {
		fmt.Println("Failed to dial to the master:", err)
		return err
	}

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	writer.WriteString(command.NewArrays([]string{"PING"}))
	writer.Flush()

	if !command.IsPong(reader) {
		return errors.New("Expected PONG from master")
	}

	writer.WriteString(command.NewArrays([]string{
		"REPLCONF",
		"listening-port",
		s.cfg.Port,
	}))
	writer.Flush()

	if !command.IsOk(reader) {
		return errors.New("Expected OK from master")
	}

	writer.WriteString(command.NewArrays([]string{
		"REPLCONF",
		"capa",
		"psync2",
	}))
	writer.Flush()

	if !command.IsOk(reader) {
		return errors.New("Expected OK from master")
	}

	writer.WriteString(command.NewArrays([]string{
		"PSYNC",
		"?",
		"-1",
	}))
	writer.Flush()

	// FULLRESYNC
	_, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	// RDB file
	_, err = reader.ReadString('\n')
	if err != nil {
		return err
	}

	go s.HandleConnection(conn)

	return nil
}

func (s *Server) Listen() (net.Listener, error) {
	return net.Listen("tcp", "0.0.0.0:"+s.cfg.Port)
}

func (s *Server) HandleConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

cmdLoop:
	for {
		redisCmd, err := command.NewRedisCommand(reader)
		if err != nil {
			fmt.Println("Unable to parse redis command:", err)
			break
		}

		fmt.Printf("Redis command: %v\n", redisCmd)

		switch strings.ToLower(redisCmd.Args[0]) {

		case command.Ping:
			writer.WriteString(command.NewSimpleString("PONG"))

		case command.Echo:
			if len(redisCmd.Args) < 2 {
				fmt.Println("Echo command requires an argument")
				break cmdLoop
			}

			echoArg := redisCmd.Args[1]

			writer.WriteString(command.NewBulkString(echoArg))

		case command.Get:
			if len(redisCmd.Args) < 2 {
				fmt.Println("Get command requires an argument")
				break cmdLoop
			}

			keyArg := redisCmd.Args[1]

			value, err := s.db.Get(keyArg)
			if err != nil {
				writer.WriteString(command.NewNulls())
				break
			}

			writer.WriteString(command.NewBulkString(value))

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

			s.db.Set(keyArg, valueArg, expArg)
			writer.WriteString(command.NewSimpleString("OK"))

			for _, conn := range s.slaves {
				fmt.Println("Try to propragate to slaves...")
				writer := bufio.NewWriter(conn)
				writer.WriteString(command.NewArrays(redisCmd.Args))
				writer.Flush()
			}

		case command.Info:
			if len(redisCmd.Args) < 2 {
				fmt.Println("Info command requires arguments")
				break cmdLoop
			}

			sectionArg := redisCmd.Args[1]

			if sectionArg == "replication" {
				info := []string{
					fmt.Sprintf("role:%s", s.cfg.Role),
					"connected_slaves:0",
					fmt.Sprintf("master_replid:%s", s.cfg.MasterReplid),
					"master_repl_offset:0",
				}
				writer.WriteString(command.NewBulkString(strings.Join(info, "\n")))
			}

		case command.Replconf:
			if len(redisCmd.Args) < 3 {
				fmt.Println("Replconf command requires arguments")
				break cmdLoop
			}

			switch redisCmd.Args[1] {

			case command.ListeningPort:
				s.slaves = append(s.slaves, conn)
				writer.WriteString(command.NewSimpleString("OK"))

			case command.GetAck:
				writer.WriteString(command.NewArrays([]string{
					"REPLCONF",
					"ACK",
					"0",
				}))

			default:
				writer.WriteString(command.NewSimpleString("OK"))
			}

		case command.Psync:
			if len(redisCmd.Args) < 3 {
				fmt.Println("Psync command requires arguments")
				break cmdLoop
			}

			resp := fmt.Sprintf("FULLRESYNC %s 0", s.cfg.MasterReplid)
			writer.WriteString(command.NewSimpleString(resp))

			emptyFile := database.EmptyRDB()
			writer.WriteString(command.NewRawBytes(emptyFile))
		}

		writer.Flush()
	}
}
