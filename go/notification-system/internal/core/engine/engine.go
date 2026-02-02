package engine

import (
	"net/http"

	"notification-system/internal/config"
	"notification-system/internal/observability/logging"
	"notification-system/internal/observability/metrics"
)

type Engine struct {
	cfg config.Config
	logger *logging.Logger
	metrics *metrics.Registry
}

func New(cfg config.Config, 
	logger *logging.Logger,
	metrics *metrics.Registry,
) *Engine {
	return &Engine{
		cfg: cfg,
		logger: logger,
		metrics: metrics,
	}
}

func (e *Engine) Health(w http.ResponseWriter, _ *http.Request) {
	e.metrics.IncHealthCheck()
	e.logger.Info("health check")
	
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

