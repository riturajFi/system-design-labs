package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type Target struct {
	Channel model.Channel
	Address string
	Device  *model.Device
	User    *model.User
}

type Resolver interface {
	Resolve(ctx context.Context, n model.Notification) ([]Target, error)
}
