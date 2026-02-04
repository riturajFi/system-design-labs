package feed

import (
	"context"
	"sync"

	"newsfeed/internal/models"
)

type MemoryStore struct {
	mu       sync.RWMutex
	maxItems int
	feeds    map[string][]models.FeedEntry
}

func NewMemoryStore(maxItems int) *MemoryStore {
	return &MemoryStore{
		maxItems: maxItems,
		feeds:    make(map[string][]models.FeedEntry),
	}
}

// Append inserts at the front (newest first).
func (m *MemoryStore) Append(_ context.Context, entry models.FeedEntry) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	list := m.feeds[entry.UserID]

	// prepend
	list = append([]models.FeedEntry{entry}, list...)

	// enforce cap
	if len(list) > m.maxItems {
		list = list[:m.maxItems]
	}

	m.feeds[entry.UserID] = list
	return nil
}

func (m *MemoryStore) Get(_ context.Context, userID string, limit int) []models.FeedEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := m.feeds[userID]

	if limit <= 0 || limit > len(list) {
		limit = len(list)
	}

	out := make([]models.FeedEntry, limit)
	copy(out, list[:limit])
	return out
}
