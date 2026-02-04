package shortening

import (
	"context"
	"fmt"

	"urlshortener/internal/domain"
	"urlshortener/internal/encoding/base62"
	"urlshortener/internal/idgen"
	"urlshortener/internal/storage/repository"
)

type Service struct {
	repo repository.URLRepository
	idg  idgen.Generator
}

func NewService(repo repository.URLRepository, idg idgen.Generator) *Service {
	return &Service{
		repo: repo,
		idg:  idg,
	}
}

// Shorten returns an existing short URL if the long URL is already known (idempotent).
// Otherwise it creates a new mapping using a newly generated ID.
func (s *Service) Shorten(ctx context.Context, longURL domain.LongURL) (domain.ShortURL, error) {
	existing, err := s.repo.GetByLongURL(ctx, longURL)
	if err != nil {
		return "", fmt.Errorf("repo get by long url: %w", err)
	}
	if existing != nil {
		return *existing, nil
	}

	id, err := s.idg.Generate(ctx)
	if err != nil {
		return "", fmt.Errorf("id generator: %w", err)
	}

	shortURL := base62.Encode(id)

	if err := s.repo.Save(ctx, id, shortURL, longURL); err != nil {
		return "", fmt.Errorf("repo save: %w", err)
	}

	return shortURL, nil
}
