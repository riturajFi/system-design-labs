package memory

import (
	"context"
	"sync"

	"urlshortener/internal/domain"
	"urlshortener/internal/idgen"
)

type Generator struct {
	mu   sync.Mutex
	next domain.ID
}

func New(startAt domain.ID) idgen.Generator {
	return &Generator{
		next: startAt,
	}
}

func (g *Generator) Generate(_ context.Context) (domain.ID, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.next++
	return g.next, nil
}
