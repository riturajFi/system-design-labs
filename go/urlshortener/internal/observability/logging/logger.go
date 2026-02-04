package logging

import (
	"log"
	"os"
)

type Logger struct {
	logger *log.Logger
}

func New(serviceName string) *Logger {
	l := log.New(
		os.Stdout,
		"",
		log.LstdFlags|log.LUTC,
	)

	l.Printf("logger initialized: service=%s", serviceName)

	return &Logger{
		logger: l,
	}
}

func (l *Logger) Info(msg string) {
	l.logger.Printf("level=INFO msg=%q", msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Printf("level=ERROR msg=%q", msg)
}
