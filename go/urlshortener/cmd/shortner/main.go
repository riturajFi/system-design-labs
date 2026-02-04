package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"urlshortener/internal/config"
	"urlshortener/internal/observability/logging"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger := logging.New(cfg.ServiceName)

	logger.Info("service starting")

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	<-ctx.Done()

	logger.Info("service shutting down")
	time.Sleep(100 * time.Millisecond) // simulate graceful cleanup
}
