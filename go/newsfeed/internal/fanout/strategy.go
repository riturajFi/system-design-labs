package fanout

import (
	"context"

	"newsfeed/internal/models"
)

// Strategy decides how a post is delivered to feeds.
type Strategy interface {
	// Distribute schedules delivery of this post.
	// Must return quickly. No blocking on fanout completion.
	Distribute(ctx context.Context, post models.Post) error
}
