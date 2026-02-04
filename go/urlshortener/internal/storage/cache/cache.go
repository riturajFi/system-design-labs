package cache

import (
	"context"

	"urlshortener/internal/domain"
)

// URLCache defines a read-optimization layer for shortURL -> longURL resolution.
// Cache is NOT a source of truth.
type URLCache interface {

	// Get returns the long URL for a short URL if present.
	// Returns (nil, nil) on cache miss.
	Get(
		ctx context.Context,
		shortURL domain.ShortURL,
	) (*domain.LongURL, error)

	// Set stores a mapping in cache.
	// Cache implementations may choose to ignore this call.
	Set(
		ctx context.Context,
		shortURL domain.ShortURL,
		longURL domain.LongURL,
	) error
}
