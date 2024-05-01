package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/handler"
)

type Server struct {
	cfg config.Config
	db  database.IDatabase
}

func NewServer(cfg config.Config, db database.IDatabase) Server {
	return Server{cfg, db}
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

	connection := handler.NewHandler(conn, s.db, s.cfg)

	if err := connection.Handshake(); err != nil {
		fmt.Println("Unable to handshake to the master:", err)
		return err
	}

	return nil
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+s.cfg.Port)
	if err != nil {
		fmt.Println("Failed to bind to port", s.cfg.Port)
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return err
		}

		connection := handler.NewHandler(conn, s.db, s.cfg)

		go connection.Handle()
	}
}
