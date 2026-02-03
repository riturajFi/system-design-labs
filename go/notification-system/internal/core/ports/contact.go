package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type ContactStore interface {
	GetUser(ctx context.Context, userID int64) (model.User, error)
	GetDevicesByUser(ctx context.Context, userID int64) ([]model.Device, error)
}
