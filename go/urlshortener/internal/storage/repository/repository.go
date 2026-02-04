package repository

import (
	"context"

	"urlshortener/internal/domain"
)

// URLRepository defines the persistence contract for URL mappings.
// It is the source of truth.
type URLRepository interface {

	// GetByShortURL returns the long URL for a given short URL.
	// Returns (nil, nil) if not found.
	GetByShortURL(
		ctx context.Context,
		shortURL domain.ShortURL,
	) (*domain.LongURL, error)

	// GetByLongURL returns the short URL for a given long URL.
	// Used to guarantee idempotency.
	// Returns (nil, nil) if not found.
	GetByLongURL(
		ctx context.Context,
		longURL domain.LongURL,
	) (*domain.ShortURL, error)

	// Save persists a new mapping.
	// Must fail if the short URL already exists.
	Save(
		ctx context.Context,
		id domain.ID,
		shortURL domain.ShortURL,
		longURL domain.LongURL,
	) error
}
