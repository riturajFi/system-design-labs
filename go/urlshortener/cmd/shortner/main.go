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
	"urlshortener/internal/idgen"
	"urlshortener/internal/idgen/memory"
	storageidgen "urlshortener/internal/idgen/storage"
	"urlshortener/internal/observability/logging"
	"urlshortener/internal/redisclient"
	"urlshortener/internal/resolution"
	"urlshortener/internal/shortening"
	"urlshortener/internal/storage/cache"
	cachememory "urlshortener/internal/storage/cache/memory"
	rediscache "urlshortener/internal/storage/cache/redis"
	"urlshortener/internal/storage/repository"
	repomemory "urlshortener/internal/storage/repository/memory"
	remoterepo "urlshortener/internal/storage/repository/remote"
	"urlshortener/internal/storageclient"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	logger := logging.New(cfg.ServiceName)
	repo, cache, idg, mode, err := buildDependencies(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
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

	logger.Info(fmt.Sprintf("service starting on %s mode=%s", server.Addr, mode))

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

func buildDependencies(cfg *config.Config) (repository.URLRepository, cache.URLCache, idgen.Generator, string, error) {
	if cfg.RedisAddr == "" && cfg.StorageBaseURL == "" {
		return repomemory.New(), cachememory.New(), memory.New(0), "memory", nil
	}

	if cfg.RedisAddr == "" || cfg.StorageBaseURL == "" {
		return nil, nil, nil, "", fmt.Errorf("REDIS_ADDR and STORAGE_BASE_URL must either both be set or both be empty")
	}

	cacheClient := redisclient.New(cfg.RedisAddr)
	storageClient := storageclient.New(cfg.StorageBaseURL)

	return remoterepo.New(storageClient), rediscache.New(cacheClient, "urlshortener:cache:"), storageidgen.New(storageClient), "redis+storage", nil
}
