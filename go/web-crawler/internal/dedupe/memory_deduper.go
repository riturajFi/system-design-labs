package dedupe

import (
	"sync"
	"web-crawler/internal/model"
)

type MemoryDeduper struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

func NewMemory() *MemoryDeduper {
	return &MemoryDeduper{
		seen: make(map[string]struct{}),
	}
}

func (d *MemoryDeduper) Seen(req model.CrawlRequest) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.seen[req.URL]
	return ok
}

func (d *MemoryDeduper) Mark(req model.CrawlRequest) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen[req.URL] = struct{}{}
}
