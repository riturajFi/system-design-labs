package worker

import (
	"context"
	"time"

	"notification-system/internal/core/ports"
	"notification-system/internal/observability/logging"
	"notification-system/internal/observability/metrics"
)

type Worker struct {
	name     string
	queue    ports.Queue
	logger   *logging.Logger
	metrics  *metrics.Registry

	templates ports.TemplateEngine
	provider  ports.Provider
	logStore  ports.LogStore
	retry     ports.RetryPolicy
}

func New(
	name string,
	queue ports.Queue,
	templates ports.TemplateEngine,
	provider ports.Provider,
	logStore ports.LogStore,
	retry ports.RetryPolicy,
	logger *logging.Logger,
	metrics *metrics.Registry,
) *Worker {
	return &Worker{
		name:      name,
		queue:     queue,
		templates: templates,
		provider:  provider,
		logStore:  logStore,
		retry:     retry,
		logger:    logger,
		metrics:   metrics,
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

			_, err = w.templates.Render(
				ctx, 
				n.TemplateID, 
				n.Channel, 
				n.Params,
			)

			if err != nil {
				continue
			}

			err = w.provider.Send(ctx, n)
			if err != nil {
				decision := w.retry.Decide(0, err)
				if decision.Retry {
					time.Sleep(decision.After)
					_ = w.queue.Enqueue(ctx, n)
				}

				continue
			}

			w.logger.Info("dequeued event " + n.EventID)
		}
	}
}
