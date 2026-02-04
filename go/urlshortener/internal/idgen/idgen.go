package idgen

import (
	"context"

	"urlshortener/internal/domain"
)

// Generator produces globally unique numeric IDs.
// Correctness and uniqueness are owned by the implementation.
type Generator interface {

	// Generate returns a new unique ID.
	// Must never return the same ID twice.
	Generate(
		ctx context.Context,
	) (domain.ID, error)
}
