package graph

import (
	"context"
	"sync"
)

type MemoryStore struct {
	mu sync.RWMutex
	friends map[string][]string
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		friends: make(map[string][]string),
	}
}

func (m *MemoryStore) AddFriendship(userA, userB string) {

	m.mu.Lock()
	defer m.mu.Unlock()
	m.friends[userA] = append(m.friends[userA], userB)
}

func (m *MemoryStore) Friends(_ context.Context, userID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	list := m.friends[userID]

	out := make([]string, len(list))
	copy(out, list)
	return out
}