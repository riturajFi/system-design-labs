package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"urlshortener/internal/config"
	"urlshortener/internal/httpapi"
	"urlshortener/internal/idgen/memory"
	"urlshortener/internal/observability/logging"
	"urlshortener/internal/resolution"
	"urlshortener/internal/shortening"
	cachememory "urlshortener/internal/storage/cache/memory"
	repomemory "urlshortener/internal/storage/repository/memory"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger := logging.New(cfg.ServiceName)
	repo := repomemory.New()
	cache := cachememory.New()
	idg := memory.New(0)
	shortener := shortening.NewService(repo, idg)
	resolver := resolution.NewService(cache, repo)
	handler := httpapi.NewHandler(shortener, resolver, cfg.BaseURL)

	mux := http.NewServeMux()
	handler.RegisterRoutes(mux)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	logger.Info(fmt.Sprintf("service starting on %s", server.Addr))

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	serverErrCh := make(chan error, 1)
	go func() {
		serverErrCh <- server.ListenAndServe()
	}()

	select {
	case err := <-serverErrCh:
		if err != nil && err != http.ErrServerClosed {
			logger.Error(fmt.Sprintf("server failed: %v", err))
			os.Exit(1)
		}
	case <-ctx.Done():
	}

	logger.Info("service shutting down")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(fmt.Sprintf("server shutdown failed: %v", err))
	}
}
