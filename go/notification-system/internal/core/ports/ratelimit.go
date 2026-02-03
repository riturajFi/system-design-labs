package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type RateLimiter interface {
	Allow(ctx context.Context, userID int64, channel model.Channel) error
}
