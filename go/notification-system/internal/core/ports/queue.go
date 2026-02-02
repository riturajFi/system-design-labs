package ports

import (
	"context"
	"notification-system/internal/core/model"
)

type Queue interface {
	Enqueue(ctx context.Context, n model.Notification) error
	Dequeue(ctx context.Context) (model.Notification, error)
	Depth(ctx context.Context) (int, error)
}