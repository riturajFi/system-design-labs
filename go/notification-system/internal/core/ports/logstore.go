package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type LogStore interface {
	Save(ctx context.Context, n model.Notification) error
	MarkSent(ctx context.Context, eventID string) error
	MarkError(ctx context.Context, eventID string, reason string) error
	Exists(ctx context.Context, eventID string) (bool, error)
}
