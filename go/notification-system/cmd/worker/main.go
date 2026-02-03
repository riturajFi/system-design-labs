package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"notification-system/internal/config"
	"notification-system/internal/observability/logging"
	"notification-system/internal/observability/metrics"
	"notification-system/internal/worker"

	queuememory "notification-system/internal/modules/queue/memory"
)

func main (){

	cfg := config.Load()
	logger := logging.New("worker", cfg.Env)
	registry := metrics.NewRegistry()

	// TODO:
	// For now, worker and API must share the SAME queue instance
	// In containers, this becomes Redis or another external queue
	queue := queuememory.New("email", registry)

	w := worker.New("email-worker", queue, logger, registry)

	ctx, cancel := context.WithCancel(context.Background())

	// TODO: Understand this necessity fo worker
	go w.Run(ctx)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	mux.Handle("/metrics", metrics.NewHandler(registry))

	go http.ListenAndServe(":8081", mux)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	// TODO: What is this?
	<- sig
	cancel()
}