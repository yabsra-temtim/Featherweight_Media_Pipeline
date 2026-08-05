package config

import (
	"os"
	"strconv"
	"time"
)

//holdes all app configuration

type Config struct {
	ServerAddress   string
	CommandTimeout  time.Duration
}

// load reads environment variables and populates the config struct
func Load() Config {
	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = ":8080" // default port
	}

	timeout := 30 * time.Second
	if raw := os.Getenv("COMMAND_TIMEOUT_SECONDS"); raw != "" {
		if seconds, err := strconv.Atoi(raw); err == nil && seconds > 0 {
			timeout = time.Duration(seconds) * time.Second
		}
	}

	return Config{
		ServerAddress:  addr,
		CommandTimeout: timeout,
	}
}
