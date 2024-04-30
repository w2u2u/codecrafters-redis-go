package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/handler"
)

type Server struct {
	port string
	db   database.IDatabase
}

func NewServer(port string, db database.IDatabase) Server {
	return Server{port, db}
}

func (Server) Slave() error {
	return nil
}

func (s *Server) Master() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+s.port)
	if err != nil {
		fmt.Println("Failed to bind to port", s.port)
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return err
		}

		connection := handler.NewHandler(conn, s.db)

		go connection.Handle()
	}
}
