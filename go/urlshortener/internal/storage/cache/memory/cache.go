package memory

import (
	"context"
	"sync"

	"urlshortener/internal/domain"
	"urlshortener/internal/storage/cache"
)

type Cache struct {
	mu   sync.RWMutex
	data map[domain.ShortURL]domain.LongURL
}

func New() cache.URLCache {
	return &Cache{
		data: make(map[domain.ShortURL]domain.LongURL),
	}
}

func (c *Cache) Get(_ context.Context, shortURL domain.ShortURL) (*domain.LongURL, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	longURL, ok := c.data[shortURL]
	if !ok {
		return nil, nil
	}

	result := longURL
	return &result, nil
}

func (c *Cache) Set(_ context.Context, shortURL domain.ShortURL, longURL domain.LongURL) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[shortURL] = longURL
	return nil
}
