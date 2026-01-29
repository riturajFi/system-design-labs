package frontier

import (
	"net/url"
	"sync"
	"web-crawler/internal/model"
)

// PriorityPoliteFrontier combines priority (front queues) with per-host politeness.
type PriorityPoliteFrontier struct {
	mu sync.Mutex

	frontQueues map[model.Priority][]model.CrawlRequest
	backQueues  map[string][]model.CrawlRequest
	inFlight    map[string]bool
}

// NewPriorityPolite constructs a frontier with priority-aware front queues and polite back queues.
func NewPriorityPolite() *PriorityPoliteFrontier {
	return &PriorityPoliteFrontier{
		frontQueues: map[model.Priority][]model.CrawlRequest{
			model.PriorityHigh: {},
			model.PriorityLow:  {},
		},
		backQueues: make(map[string][]model.CrawlRequest),
		inFlight:   make(map[string]bool),
	}
}

// Push places a request into its priority-specific front queue.
func (p *PriorityPoliteFrontier) Push(req model.CrawlRequest) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.frontQueues[req.Priority] = append(p.frontQueues[req.Priority], req)
}

// Pop selects from front queues (priority-first) and enforces per-host politeness.
func (p *PriorityPoliteFrontier) Pop() (model.CrawlRequest, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Biased selection: try high priority before low priority.
	for _, pr := range []model.Priority{model.PriorityHigh, model.PriorityLow} {
		if len(p.frontQueues[pr]) == 0 {
			continue
		}

		req := p.frontQueues[pr][0]
		p.frontQueues[pr] = p.frontQueues[pr][1:]

		u, err := url.Parse(req.URL)
		if err != nil || u.Host == "" {
			continue
		}

		if p.inFlight[u.Host] {
			// Host is busy, push back to same priority queue.
			p.frontQueues[pr] = append(p.frontQueues[pr], req)
			continue
		}

		// Move into the per-host back queue and mark in-flight.
		p.inFlight[u.Host] = true
		p.backQueues[u.Host] = append(p.backQueues[u.Host], req)
		next := p.backQueues[u.Host][0]
		p.backQueues[u.Host] = p.backQueues[u.Host][1:]
		return next, true
	}

	return model.CrawlRequest{}, false
}

// Done marks a host as no longer in-flight after a request finishes.
func (p *PriorityPoliteFrontier) Done(req model.CrawlRequest) {
	u, err := url.Parse(req.URL)
	if err != nil || u.Host == "" {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.inFlight[u.Host] = false
}

// Len returns the total number of queued requests across all priority levels.
func (p *PriorityPoliteFrontier) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	total := 0
	for _, q := range p.frontQueues {
		total += len(q)
	}
	return total
}
