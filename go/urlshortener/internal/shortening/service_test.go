package shortening

import (
	"context"
	"sync"
	"testing"

	"urlshortener/internal/domain"
	"urlshortener/internal/idgen/memory"
	repomemory "urlshortener/internal/storage/repository/memory"
)

func TestShortenIsIdempotentUnderConcurrency(t *testing.T) {
	repo := repomemory.New()
	idg := memory.New(0)
	service := NewService(repo, idg)

	const workers = 16

	results := make([]domain.ShortURL, workers)
	errs := make([]error, workers)

	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			results[i], errs[i] = service.Shorten(context.Background(), domain.LongURL("https://example.com/concurrent"))
		}(i)
	}

	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("worker %d returned error: %v", i, err)
		}
		if results[i] != "1" {
			t.Fatalf("worker %d returned short code %q, want %q", i, results[i], "1")
		}
	}
}
