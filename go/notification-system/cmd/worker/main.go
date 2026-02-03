package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"notification-system/internal/config"
	"notification-system/internal/core/model"
	logmemory "notification-system/internal/modules/logstore/memory"
	provideremail "notification-system/internal/modules/providers/email/mock"
	queuememory "notification-system/internal/modules/queue/memory"
	retryfixed "notification-system/internal/modules/retry/fixed"
	templatememory "notification-system/internal/modules/template/memory"
	httptracker "notification-system/internal/modules/tracking/http"
	"notification-system/internal/observability/logging"
	"notification-system/internal/observability/metrics"
	"notification-system/internal/worker"
)

func main (){

	cfg := config.Load()
	logger := logging.New("worker", cfg.Env)
	registry := metrics.NewRegistry()

	// TODO:
	// For now, worker and API must share the SAME queue instance
	// In containers, this becomes Redis or another external queue
	queue := queuememory.New("email", registry)
	templates := templatememory.New(map[string]map[model.Channel]templatememory.Template{
		"welcome": {
			model.ChannelEmail: {
				Title: "Welcome {{name}}",
				Body:  "Hello {{name}}, thanks for joining!",
			},
		},
	})

	provider := provideremail.New(true) // no failure
	logStore := logmemory.New()
	retryPolicy := retryfixed.New(3, time.Second)
	tracker := httptracker.New("http://localhost:8090/ingest")

	w := worker.New(
		"email-worker",
		queue,
		templates,
		provider,
		logStore,
		retryPolicy,
		tracker,
		logger,
		registry,
	)

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
