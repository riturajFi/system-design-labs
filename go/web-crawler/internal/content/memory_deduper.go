package content

import "sync"

// MemoryDeduper stores seen content hashes in memory.
type MemoryDeduper struct {
	mu   sync.Mutex
	seen map[string]struct{}
}

// NewMemory constructs an in-memory content deduper.
func NewMemory() *MemoryDeduper {
	return &MemoryDeduper{
		seen: make(map[string]struct{}),
	}
}

// Seen checks whether a content hash has been observed before.
func (d *MemoryDeduper) Seen(hash string) bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	_, ok := d.seen[hash]
	return ok
}

// Mark records a content hash as seen.
func (d *MemoryDeduper) Mark(hash string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.seen[hash] = struct{}{}
}
