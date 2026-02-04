package config

import (
	"fmt"
	"os"
)

type Config struct {
	ServiceName string
	Port        string
}

func Load() (*Config, error) {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		return nil, fmt.Errorf("SERVICE_NAME is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		return nil, fmt.Errorf("PORT is required")
	}

	return &Config{
		ServiceName: serviceName,
		Port:        port,
	}, nil
}
