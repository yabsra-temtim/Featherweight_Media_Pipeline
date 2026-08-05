package config

import "os"

//holdes all app configuration

type Config struct {
	ServerAddress string
}

// load reads environment variables and populates the config struct
func Load() Config {
	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = ":8080" //default port
	}
	return Config{
		ServerAddress: addr,
	}
}
