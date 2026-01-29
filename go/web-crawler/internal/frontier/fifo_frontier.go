package frontier

import (
	"sync"
	"web-crawler/internal/model"
)

type FIFOFrontier struct {
	mu    sync.Mutex
	queue []model.CrawlRequest
}

func NewFIFO() *FIFOFrontier {
	return &FIFOFrontier{
		queue: make([]model.CrawlRequest, 0),
	}
}

func (f *FIFOFrontier) Push(req model.CrawlRequest) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.queue = append(f.queue, req)
}

func (f *FIFOFrontier) Pop() (model.CrawlRequest, bool) {
	f.mu.Lock()
	defer f.mu.Unlock()

	if len(f.queue) == 0 {
		return model.CrawlRequest{}, false
	}

	req := f.queue[0]
	f.queue = f.queue[1:]
	return req, true
}

func (f *FIFOFrontier) Len() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.queue)
}
