package fixed

import (
	"context"
	"errors"
	"sync"
	"time"

	"notification-system/internal/core/model"
	"notification-system/internal/observability/metrics"
)

type key struct {
	userID  int64
	channel model.Channel
	window  int64
}

type FixedWindowLimiter struct {
	mu      sync.Mutex
	limit   int
	windows map[key]int
	metrics *metrics.Registry
}

func New(limit int, metrics *metrics.Registry) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		limit:   limit,
		windows: make(map[key]int),
		metrics: metrics,
	}
}

func (l *FixedWindowLimiter) Allow(
	_ context.Context,
	userID int64,
	channel model.Channel,
) error {
	nowWindow := time.Now().Unix() / 60

	k := key{
		userID:  userID,
		channel: channel,
		window:  nowWindow,
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	count := l.windows[k]
	if count >= l.limit {
		if l.metrics != nil {
			l.metrics.IncRateDenied()
		}
		return errors.New("rate limit exceeded")
	}

	l.windows[k] = count + 1
	if l.metrics != nil {
		l.metrics.IncRateAllowed()
	}
	return nil
}
