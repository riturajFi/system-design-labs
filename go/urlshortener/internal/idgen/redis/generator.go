package redis

import (
	"context"
	"fmt"

	"urlshortener/internal/domain"
	"urlshortener/internal/idgen"
	"urlshortener/internal/redisclient"
)

type Generator struct {
	client *redisclient.Client
	key    string
}

func New(client *redisclient.Client, key string) idgen.Generator {
	return &Generator{
		client: client,
		key:    key,
	}
}

func (g *Generator) Generate(ctx context.Context) (domain.ID, error) {
	id, err := g.client.Incr(ctx, g.key)
	if err != nil {
		return 0, fmt.Errorf("redis incr: %w", err)
	}

	return domain.ID(id), nil
}
