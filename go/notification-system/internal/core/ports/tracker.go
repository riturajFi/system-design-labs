package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type Tracker interface {
	Track(ctx context.Context, ev model.NotificationEvent) error
}
