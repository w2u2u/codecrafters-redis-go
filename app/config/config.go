package config

import (
	"fmt"
	"net"
	"os"
)

type Config struct {
	Port         string
	ReplicaOf    string
	Role         string
	MasterReplid string
	Slaves       []net.Conn
}

func Parse() Config {
	cfg := Config{
		Port:         "6379",
		Role:         "master",
		MasterReplid: "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb",
	}
	argsWithoutProg := os.Args[1:]

	for i, arg := range argsWithoutProg {
		if arg == "--port" && len(argsWithoutProg) > i {
			cfg.Port = argsWithoutProg[i+1]
		}

		if arg == "--replicaof" && len(argsWithoutProg) > i {
			cfg.Role = "slave"
			cfg.ReplicaOf = fmt.Sprintf("%s:%s", argsWithoutProg[i+1], argsWithoutProg[i+2])
		}
	}

	return cfg
}

func (cfg *Config) AddSlave(conn net.Conn) {
	cfg.Slaves = append(cfg.Slaves, conn)
}
