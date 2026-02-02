package engine

import (
	"net/http"

	"notification-system/internal/config"
	"notification-system/internal/observability/logging"
)

type Engine struct {
	cfg config.Config
	logger *logging.Logger
}

func New(cfg config.Config, logger *logging.Logger) *Engine {
	return &Engine{
		cfg: cfg,
		logger: logger,
	}
}

func (e *Engine) Health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

