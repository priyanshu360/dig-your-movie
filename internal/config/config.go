package config

import (
	"os"
	"strconv"
)

type Config struct {
	APIKey string
	Port   int
}

func Load() *Config {
	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		apiKey = "9496f5e1" // Default for dev
	}

	portStr := os.Getenv("PORT")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8095 // Default port
	}

	return &Config{
		APIKey: apiKey,
		Port:   port,
	}
}
