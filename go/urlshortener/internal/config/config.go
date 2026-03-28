package config

import (
	"os"
	"strings"
)

type Config struct {
	ServiceName string
	Port        string
	BaseURL     string
}

func Load() (*Config, error) {
	serviceName := os.Getenv("SERVICE_NAME")
	if serviceName == "" {
		serviceName = "urlshortener"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		ServiceName: serviceName,
		Port:        port,
		BaseURL:     strings.TrimRight(strings.TrimSpace(os.Getenv("BASE_URL")), "/"),
	}, nil
}
