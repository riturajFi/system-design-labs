package ratelimit

import (
	"context"
	"sync"
	"time"
)

type MemoryLimiter struct {
	mu     sync.Mutex
	limit  int
	window time.Duration
	hits   map[string][]time.Time
}

func NewMemoryLimiter(limit int, window time.Duration) *MemoryLimiter {
	return &MemoryLimiter{
		limit:  limit,
		window: window,
		hits:   make(map[string][]time.Time),
	}
}

func (m *MemoryLimiter) Allow(_ context.Context, key string) bool {
	now := time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()

	events := m.hits[key]
	valid := events[:0]

	for _, t := range events {
		if now.Sub(t) < m.window {
			valid = append(valid, t)
		}
	}

	if len(valid) >= m.limit {
		m.hits[key] = valid
		return false
	}

	valid = append(valid, now)
	m.hits[key] = valid
	return true
}
