package worker

import (
	"context"
	"time"

	"notification-system/internal/core/model"
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
	tracker   ports.Tracker
}

func New(
	name string,
	queue ports.Queue,
	templates ports.TemplateEngine,
	provider ports.Provider,
	logStore ports.LogStore,
	retry ports.RetryPolicy,
	tracker ports.Tracker,
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
		tracker:   tracker,
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

			// Defensive deduplication (worker-side)
			exists, err := w.logStore.Exists(ctx, n.EventID)
			if err != nil {
				w.logger.Error("logstore error for event " + n.EventID)
				continue
			}

			if exists {
				// Already processed or in-progress
				w.logger.Info("duplicate event skipped: " + n.EventID)
				continue
			}

			// Mark intent before delivery to avoid races
			if err := w.logStore.Save(ctx, n); err != nil {
				w.logger.Info("duplicate detected on save: " + n.EventID)
				continue
			}

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
				} else {
					trackErr := w.tracker.Track(ctx, model.NotificationEvent{
						EventID: n.EventID,
						UserID:  n.UserID,
						Channel: n.Channel,
						Type:    model.EventError,
						Message: err.Error(),
						AtUnix:  time.Now().Unix(),
					})
					if trackErr != nil {
						w.metrics.IncTrackingFailed()
					}
				}

				continue
			}

			trackErr := w.tracker.Track(ctx, model.NotificationEvent{
				EventID: n.EventID,
				UserID:  n.UserID,
				Channel: n.Channel,
				Type:    model.EventSent,
				Message: "sent",
				AtUnix:  time.Now().Unix(),
			})
			if trackErr != nil {
				w.metrics.IncTrackingFailed()
			}

			w.logger.Info("dequeued event " + n.EventID)
		}
	}
}
