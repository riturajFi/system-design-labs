package memory

import (
	"context"
	"errors"
	"sync"

	"notification-system/internal/core/model"
	"notification-system/internal/observability/metrics"
)

type Queue struct {
	name    string
	mu      sync.Mutex
	items   []model.Notification
	metrics *metrics.Registry
}

func New(name string, metrics *metrics.Registry) *Queue {
	return &Queue{
		name:    name,
		metrics: metrics,
	}
}

func (q *Queue) Enqueue(
	_ context.Context,
	n model.Notification,
) error {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, n)
	if q.metrics != nil {
		q.metrics.SetQueueDepth(q.name, len(q.items))
	}
	return nil
}

func (q *Queue) Dequeue(
	_ context.Context,
) (model.Notification, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) == 0 {
		return model.Notification{}, errors.New("queue empty")
	}

	n := q.items[0]
	q.items = q.items[1:]
	if q.metrics != nil {
		q.metrics.SetQueueDepth(q.name, len(q.items))
	}
	return n, nil
}

func (q *Queue) Depth(
	_ context.Context,
) (int, error) {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items), nil
}
