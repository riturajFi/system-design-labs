package dedupe

import "web-crawler/internal/model"

type MemoryDeduper struct {
	seen map[string]struct{}
}

func NewMemory() *MemoryDeduper {
	return &MemoryDeduper{
		seen: make(map[string]struct{}),
	}
}

func (d *MemoryDeduper) Seen(req model.CrawlRequest) bool {
	_, ok := d.seen[req.URL]
	return ok
}

func (d *MemoryDeduper) Mark(req model.CrawlRequest) {
	d.seen[req.URL] = struct{}{}
}
