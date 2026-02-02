package config

import (
	"os"
)

type Config struct {
	ServiceName string
	Env         string
	HTTPPort    string
}

func Load() Config {
	return Config{
		ServiceName: getEnv("SERVICE_NAME", "notification-api"),
		Env:         getEnv("ENV", "dev"),
		HTTPPort:    getEnv("HTTP_PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
