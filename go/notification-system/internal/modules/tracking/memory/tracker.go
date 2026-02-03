package memory

import (
	"context"
	"sync"

	"notification-system/internal/core/model"
	"notification-system/internal/observability/metrics"
)

type Tracker struct {
	mu      sync.Mutex
	events  []model.NotificationEvent
	metrics *metrics.Registry
}

func New(metrics *metrics.Registry) *Tracker {
	return &Tracker{
		metrics: metrics,
	}
}

func (t *Tracker) Track(
	_ context.Context,
	e model.NotificationEvent,
) error {

	t.mu.Lock()
	t.events = append(t.events, e)
	t.mu.Unlock()

	t.metrics.IncTrackingEmitted()
	return nil
}
