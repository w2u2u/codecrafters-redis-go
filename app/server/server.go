package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/database"
	"github.com/codecrafters-io/redis-starter-go/app/handler"
)

type Server struct{ port string }

func NewServer(port string) Server {
	return Server{port}
}

func (Server) Slave() error {
	return nil
}

func (s Server) Master() error {
	l, err := net.Listen("tcp", "0.0.0.0:"+s.port)
	if err != nil {
		fmt.Println("Failed to bind to port", s.port)
		return err
	}

	keyValueDb := database.NewKeyValue()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			return err
		}

		connection := handler.NewHandler(conn, &keyValueDb)

		go connection.Handle()
	}
}
