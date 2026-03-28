package storage

import (
	"context"
	"fmt"

	"urlshortener/internal/domain"
	"urlshortener/internal/idgen"
	"urlshortener/internal/storageclient"
)

type Generator struct {
	client *storageclient.Client
}

func New(client *storageclient.Client) idgen.Generator {
	return &Generator{client: client}
}

func (g *Generator) Generate(ctx context.Context) (domain.ID, error) {
	id, err := g.client.NextID(ctx)
	if err != nil {
		return 0, fmt.Errorf("storage next id: %w", err)
	}

	return id, nil
}
