package engine

import (
	"net/http"

	"notification-system/internal/config"
	"notification-system/internal/core/ports"
	"notification-system/internal/observability/logging"
	"notification-system/internal/observability/metrics"
)

type Deps struct {
	Auth      ports.Authenticator
	RateLimit ports.RateLimiter
	Settings  ports.SettingsChecker
	Queue     ports.Queue
	LogStore  ports.LogStore
	Tracker   ports.Tracker
}

type Engine struct {
	cfg     config.Config
	logger  *logging.Logger
	metrics *metrics.Registry
	deps    Deps
}

func New(
	cfg config.Config,
	logger *logging.Logger,
	metrics *metrics.Registry,
	deps Deps,
) *Engine {
	return &Engine{
		cfg:     cfg,
		logger:  logger,
		metrics: metrics,
		deps:    deps,
	}
}

func (e *Engine) Health(w http.ResponseWriter, _ *http.Request) {
	e.metrics.IncHealthCheck()
	e.logger.Info("health check")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
