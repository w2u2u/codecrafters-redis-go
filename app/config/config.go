package config

import "os"

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
			cfg.ReplicaOf = argsWithoutProg[i+1]
		}
	}

	return cfg
}
