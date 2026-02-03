package static

import (
	"context"

	"notification-system/internal/core/model"
)

type Key struct {
	UserID  int64
	Channel model.Channel
}

type StaticSettings struct {
	optIn map[Key]bool
}

func New(optIn map[Key]bool) *StaticSettings {
	return &StaticSettings{optIn: optIn}
}

func (s *StaticSettings) IsOptedIn(
	_ context.Context,
	userID int64,
	channel model.Channel,
) (bool, error) {

	v, ok := s.optIn[Key{
		UserID:  userID,
		Channel: channel,
	}]
	if !ok {
		// default: opted out if not explicitly present
		return false, nil
	}
	return v, nil
}
