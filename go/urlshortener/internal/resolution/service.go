package resolution

import (
	"context"
	"fmt"

	"urlshortener/internal/domain"
	"urlshortener/internal/storage/cache"
	"urlshortener/internal/storage/repository"
)

type Service struct {
	cache cache.URLCache
	repo  repository.URLRepository
}

func NewService(cache cache.URLCache, repo repository.URLRepository) *Service {
	return &Service{
		cache: cache,
		repo:  repo,
	}
}

func (s *Service) Resolve(ctx context.Context, shortURL domain.ShortURL) (domain.LongURL, error) {
	if s.cache != nil {
		cached, err := s.cache.Get(ctx, shortURL)
		if err != nil {
			return "", fmt.Errorf("cache get: %w", err)
		}
		if cached != nil {
			return *cached, nil
		}
	}

	longURL, err := s.repo.GetByShortURL(ctx, shortURL)
	if err != nil {
		return "", fmt.Errorf("repo get by short url: %w", err)
	}
	if longURL == nil {
		return "", nil
	}

	if s.cache != nil {
		if err := s.cache.Set(ctx, shortURL, *longURL); err != nil {
			return "", fmt.Errorf("cache set: %w", err)
		}
	}

	return *longURL, nil
}
