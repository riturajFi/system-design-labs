package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"urlshortener/internal/domain"
	"urlshortener/internal/redisclient"
	"urlshortener/internal/storage/repository"
)

type Repository struct {
	client *redisclient.Client
}

func New(client *redisclient.Client) repository.URLRepository {
	return &Repository{client: client}
}

func (r *Repository) GetByShortURL(ctx context.Context, shortURL domain.ShortURL) (*domain.LongURL, error) {
	longURL, ok, err := r.client.Get(ctx, shortKey(shortURL))
	if err != nil {
		return nil, fmt.Errorf("redis get short url: %w", err)
	}
	if !ok {
		return nil, nil
	}

	result := domain.LongURL(longURL)
	return &result, nil
}

func (r *Repository) GetByLongURL(ctx context.Context, longURL domain.LongURL) (*domain.ShortURL, error) {
	shortURL, ok, err := r.client.Get(ctx, longKey(longURL))
	if err != nil {
		return nil, fmt.Errorf("redis get long url: %w", err)
	}
	if !ok {
		return nil, nil
	}

	result := domain.ShortURL(shortURL)
	return &result, nil
}

func (r *Repository) Save(ctx context.Context, _ domain.ID, shortURL domain.ShortURL, longURL domain.LongURL) error {
	longSaved, err := r.client.SetNX(ctx, longKey(longURL), string(shortURL))
	if err != nil {
		return fmt.Errorf("redis save long mapping: %w", err)
	}
	if !longSaved {
		return fmt.Errorf("long url already exists: %s", longURL)
	}

	shortSaved, err := r.client.SetNX(ctx, shortKey(shortURL), string(longURL))
	if err != nil {
		_ = r.client.Del(ctx, longKey(longURL))
		return fmt.Errorf("redis save short mapping: %w", err)
	}
	if !shortSaved {
		_ = r.client.Del(ctx, longKey(longURL))
		return fmt.Errorf("short url already exists: %s", shortURL)
	}

	return nil
}

func shortKey(shortURL domain.ShortURL) string {
	return "urlshortener:short:" + string(shortURL)
}

func longKey(longURL domain.LongURL) string {
	sum := sha256.Sum256([]byte(longURL))
	return "urlshortener:long:" + hex.EncodeToString(sum[:])
}
