package redis

import (
	"context"
	"fmt"

	"urlshortener/internal/domain"
	"urlshortener/internal/redisclient"
	"urlshortener/internal/storage/cache"
)

type Cache struct {
	client *redisclient.Client
	prefix string
}

func New(client *redisclient.Client, prefix string) cache.URLCache {
	return &Cache{
		client: client,
		prefix: prefix,
	}
}

func (c *Cache) Get(ctx context.Context, shortURL domain.ShortURL) (*domain.LongURL, error) {
	longURL, ok, err := c.client.Get(ctx, c.key(shortURL))
	if err != nil {
		return nil, fmt.Errorf("redis cache get: %w", err)
	}
	if !ok {
		return nil, nil
	}

	result := domain.LongURL(longURL)
	return &result, nil
}

func (c *Cache) Set(ctx context.Context, shortURL domain.ShortURL, longURL domain.LongURL) error {
	if err := c.client.Set(ctx, c.key(shortURL), string(longURL)); err != nil {
		return fmt.Errorf("redis cache set: %w", err)
	}

	return nil
}

func (c *Cache) key(shortURL domain.ShortURL) string {
	return c.prefix + string(shortURL)
}
