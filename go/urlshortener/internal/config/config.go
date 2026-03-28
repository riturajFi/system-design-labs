package config

import (
	"os"
	"strings"
)

type Config struct {
	ServiceName    string
	Port           string
	BaseURL        string
	RedisAddr      string
	StorageBaseURL string
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
		ServiceName:    serviceName,
		Port:           port,
		BaseURL:        strings.TrimRight(strings.TrimSpace(os.Getenv("BASE_URL")), "/"),
		RedisAddr:      strings.TrimSpace(os.Getenv("REDIS_ADDR")),
		StorageBaseURL: strings.TrimRight(strings.TrimSpace(os.Getenv("STORAGE_BASE_URL")), "/"),
	}, nil
}
