package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"urlshortener/internal/observability/logging"
	filerepo "urlshortener/internal/storage/repository/file"
	"urlshortener/internal/storagehttp"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	filePath := os.Getenv("STORAGE_FILE")
	if filePath == "" {
		filePath = "/data/mappings.json"
	}

	repo, err := filerepo.Open(filePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger := logging.New("urlshortener-storage")
	handler := storagehttp.NewHandler(repo, repo)
	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info(fmt.Sprintf("storage service starting on %s file=%s", server.Addr, filePath))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrCh:
		if err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("storage server failed: %v", err))
			os.Exit(1)
		}
	case <-ctx.Done():
	}

	logger.Info("storage service shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("storage shutdown failed: %v", err))
	}
}
