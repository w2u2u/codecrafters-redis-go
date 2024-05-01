package config

import (
	"fmt"
	"os"
)

type Config struct {
	Port      string
	ReplicaOf string
	Role      string
}

func Parse() Config {
	cfg := Config{Port: "6379", Role: "master"}
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
