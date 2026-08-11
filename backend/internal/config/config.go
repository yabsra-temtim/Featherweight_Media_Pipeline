package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	UploadDirectory string
	ServerAddress   string
	CommandTimeout  time.Duration
}

func Load() Config {
	addr := os.Getenv("SERVER_ADDRESS")
	if addr == "" {
		addr = ":8080"
	}

	uploadDir := os.Getenv("UPLOAD_DIRECTORY")
	if uploadDir == "" {
		uploadDir = "uploads/originals"
	}

	timeout := 30 * time.Second

	if raw := os.Getenv("COMMAND_TIMEOUT_SECONDS"); raw != "" {
		if seconds, err := strconv.Atoi(raw); err == nil && seconds > 0 {
			timeout = time.Duration(seconds) * time.Second
		}
	}

	return Config{
		UploadDirectory: uploadDir,
		ServerAddress:   addr,
		CommandTimeout:  timeout,
	}
}
