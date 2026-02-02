package ports

import (
	"context"

	"notification-system/internal/core/model"
)

type SettingsChecker interface {
	IsOptedIn(ctx context.Context, userID string, channel model.Channel) (bool, error)
}
