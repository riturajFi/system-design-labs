package main

import (
	"net/http"

	"notification-system/internal/config"
	"notification-system/internal/core/engine"
	"notification-system/internal/observability/logging"
)

func main() {
	cfg := config.Load()
	logger := logging.New(cfg.ServiceName, cfg.Env)

	e := engine.New(cfg, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", e.Health)

	logger.Info("starting http server on :" + cfg.HTTPPort)
	http.ListenAndServe(":"+cfg.HTTPPort, mux)
}
