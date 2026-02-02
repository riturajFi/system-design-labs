package engine

import (
	"context"
	"errors"

	"notification-system/internal/core/model"
)

func (e *Engine) Send(
	ctx context.Context,
	appKey string,
	appSecret string,
	n model.Notification,
) error {
	if e.deps.Auth == nil ||
		e.deps.RateLimit == nil ||
		e.deps.Settings == nil ||
		e.deps.Queue == nil ||
		e.deps.LogStore == nil ||
		e.deps.Tracker == nil {
		return errors.New("engine dependencies not configured")
	}

	if err := e.deps.Auth.Authenticate(ctx, appKey, appSecret); err != nil {
		return err
	}

	if err := e.deps.RateLimit.Allow(ctx, n.UserID, n.Channel); err != nil {
		return err
	}

	ok, err := e.deps.Settings.IsOptedIn(ctx, n.UserID, n.Channel)
	if err != nil {
		return err
	}
	if !ok {
		return errors.New("user opted out")
	}

	exists, err := e.deps.LogStore.Exists(ctx, n.EventID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("duplicate event")
	}

	if err := e.deps.LogStore.Save(ctx, n); err != nil {
		return err
	}

	if err := e.deps.Queue.Enqueue(ctx, n); err != nil {
		return err
	}

	_ = e.deps.Tracker.Track(ctx, model.NotificationEvent{
		EventID: n.EventID,
		UserID:  n.UserID,
		Channel: n.Channel,
		Type:    model.EventPending,
	})

	return nil
}
