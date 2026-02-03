package worker

import (
	"context"
	"time"

	"notification-system/internal/core/ports"
	"notification-system/internal/observability/logging"
	"notification-system/internal/observability/metrics"
)

type Worker struct {
	name    string
	queue   ports.Queue
	logger  *logging.Logger
	metrics *metrics.Registry
}

func New(
	name string,
	queue ports.Queue,
	logger *logging.Logger,
	metrics *metrics.Registry,
) *Worker {
	return &Worker{
		name:    name,
		queue:   queue,
		logger:  logger,
		metrics: metrics,
	}
}

func (w *Worker) Run(ctx context.Context) {
	w.logger.Info("worker started: " + w.name)

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("worker stopped: " + w.name)
			return
		default:
			n, err := w.queue.Dequeue(ctx)
			if err != nil {
				time.Sleep(500 * time.Millisecond)
				continue
			}

			w.metrics.IncWorkerDequeued()
			w.logger.Info("dequeued event " + n.EventID)
		}
	}
}
