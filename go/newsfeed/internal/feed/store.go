package feed

import (
	"context"

	"newsfeed/internal/models"
)

// Store owns feed persistence.
// Feed is ordered newest -> oldest.
type Store interface {
	Append(ctx context.Context, entry models.FeedEntry) error
	Get(ctx context.Context, userID string, limit int) []models.FeedEntry
}
