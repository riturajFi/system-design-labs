package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type Queue interface {
	Enqueue(ctx context.Context, n model.Notification) error
}
