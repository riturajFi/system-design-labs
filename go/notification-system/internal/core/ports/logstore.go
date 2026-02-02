package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type LogStore interface {
	Exists(ctx context.Context, eventID string) (bool, error)
	Save(ctx context.Context, n model.Notification) error
}
