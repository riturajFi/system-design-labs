package logging

import (
	"log"
	"os"
)

type Logger struct {
	service string
	env string
	base *log.Logger 
}

func New(service, env string) *Logger {

	return  &Logger{
		service: service,
		env: env,
		base: log.New(os.Stdout, "", log.LstdFlags),
	}
}


func (l *Logger) Info(msg string) {
	l.base.Printf(
		`{"level":"info","service":"%s","env":"%s","msg":"%s"}`,
		l.service,
		l.env,
		msg,
	)
}

func (l *Logger) Error(msg string) {
	l.base.Printf(
		`{"level":"error","service":"%s","env":"%s","msg":"%s"}`,
		l.service,
		l.env,
		msg,
	)
}