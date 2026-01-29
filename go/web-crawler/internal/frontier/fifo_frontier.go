package frontier

import "web-crawler/internal/model"

type FIFOFrontier struct {
	queue []model.CrawlRequest
}

func NewFIFO() *FIFOFrontier {
	return &FIFOFrontier{
		queue: make([]model.CrawlRequest, 0),
	}
}

func (f *FIFOFrontier) Push(req model.CrawlRequest) {
	f.queue = append(f.queue, req)
}

func (f *FIFOFrontier) Pop() (model.CrawlRequest, bool) {
	if len(f.queue) == 0 {
		return model.CrawlRequest{}, false
	}

	req := f.queue[0]
	f.queue = f.queue[1:]
	return req, true
}

func (f *FIFOFrontier) Len() int {
	return len(f.queue)
}
