package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type Provider interface {
	Send(ctx context.Context, n model.Notification) error
}
