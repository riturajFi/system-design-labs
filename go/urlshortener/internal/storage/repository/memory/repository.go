package memory

import (
	"context"
	"fmt"
	"sync"

	"urlshortener/internal/domain"
	"urlshortener/internal/storage/repository"
)

type Repository struct {
	mu         sync.RWMutex
	byShortURL map[domain.ShortURL]domain.LongURL
	byLongURL  map[domain.LongURL]domain.ShortURL
}

func New() repository.URLRepository {
	return &Repository{
		byShortURL: make(map[domain.ShortURL]domain.LongURL),
		byLongURL:  make(map[domain.LongURL]domain.ShortURL),
	}
}

func (r *Repository) GetByShortURL(_ context.Context, shortURL domain.ShortURL) (*domain.LongURL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	longURL, ok := r.byShortURL[shortURL]
	if !ok {
		return nil, nil
	}

	result := longURL
	return &result, nil
}

func (r *Repository) GetByLongURL(_ context.Context, longURL domain.LongURL) (*domain.ShortURL, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	shortURL, ok := r.byLongURL[longURL]
	if !ok {
		return nil, nil
	}

	result := shortURL
	return &result, nil
}

func (r *Repository) Save(_ context.Context, id domain.ID, shortURL domain.ShortURL, longURL domain.LongURL) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byShortURL[shortURL]; exists {
		return fmt.Errorf("short url already exists: %s", shortURL)
	}
	if _, exists := r.byLongURL[longURL]; exists {
		return fmt.Errorf("long url already exists: %s", longURL)
	}

	r.byShortURL[shortURL] = longURL
	r.byLongURL[longURL] = shortURL

	return nil
}
