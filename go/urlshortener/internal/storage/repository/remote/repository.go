package remote

import (
	"context"
	"fmt"

	"urlshortener/internal/domain"
	"urlshortener/internal/storage/repository"
	"urlshortener/internal/storageclient"
)

type Repository struct {
	client *storageclient.Client
}

func New(client *storageclient.Client) repository.URLRepository {
	return &Repository{client: client}
}

func (r *Repository) GetByShortURL(ctx context.Context, shortURL domain.ShortURL) (*domain.LongURL, error) {
	mapping, err := r.client.GetByShortURL(ctx, shortURL)
	if err != nil {
		return nil, fmt.Errorf("storage get by short: %w", err)
	}
	if mapping == nil {
		return nil, nil
	}

	result := mapping.LongURL
	return &result, nil
}

func (r *Repository) GetByLongURL(ctx context.Context, longURL domain.LongURL) (*domain.ShortURL, error) {
	mapping, err := r.client.GetByLongURL(ctx, longURL)
	if err != nil {
		return nil, fmt.Errorf("storage get by long: %w", err)
	}
	if mapping == nil {
		return nil, nil
	}

	result := mapping.ShortURL
	return &result, nil
}

func (r *Repository) Save(ctx context.Context, id domain.ID, shortURL domain.ShortURL, longURL domain.LongURL) error {
	if err := r.client.Save(ctx, storageclient.Mapping{
		ID:       id,
		ShortURL: shortURL,
		LongURL:  longURL,
	}); err != nil {
		return fmt.Errorf("storage save: %w", err)
	}

	return nil
}
