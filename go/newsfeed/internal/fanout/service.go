package fanout

import (
	"context"

	"newsfeed/internal/models"
)

type Service struct {
	strategy Strategy
}

func NewService(strategy Strategy) *Service {
	return &Service{strategy: strategy}
}

func (s *Service) Fanout(ctx context.Context, post models.Post) error {
	return s.strategy.Distribute(ctx, post)
}
