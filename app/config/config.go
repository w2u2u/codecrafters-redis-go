package config

import "os"

type config struct {
	Port string
}

func Parse() config {
	cfg := config{Port: "6379"}
	argsWithoutProg := os.Args[1:]

	for i, arg := range argsWithoutProg {
		if arg == "--port" && len(argsWithoutProg) > i {
			cfg.Port = argsWithoutProg[i+1]
		}
	}

	return cfg
}
