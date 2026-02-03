package memory

import (
	"context"
	"errors"
	"sync"

	"notification-system/internal/core/model"
)

type entry struct {
	notification model.Notification
	state        string
	errorMsg     string
}

type Store struct {
	mu     sync.Mutex
	events map[string]*entry
}

func New() *Store {
	return &Store{
		events: make(map[string]*entry),
	}
}

func (s *Store) Exists(
	_ context.Context,
	eventID string,
) (bool, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.events[eventID]
	return ok, nil
}

func (s *Store) Save(
	_ context.Context,
	n model.Notification,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.events[n.EventID]; exists {
		return errors.New("event already exists")
	}

	s.events[n.EventID] = &entry{
		notification: n,
		state:        "pending",
	}
	return nil
}

func (s *Store) MarkSent(
	_ context.Context,
	eventID string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.events[eventID]
	if !ok {
		return errors.New("event not found")
	}

	e.state = "sent"
	return nil
}

func (s *Store) MarkError(
	_ context.Context,
	eventID string,
	reason string,
) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.events[eventID]
	if !ok {
		return errors.New("event not found")
	}

	e.state = "error"
	e.errorMsg = reason
	return nil
}
