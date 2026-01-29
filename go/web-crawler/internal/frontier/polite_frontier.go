package frontier

import (
	"net/url"
	"sync"
	"web-crawler/internal/model"
)

type PoliteFrontier struct {
	mu        sync.Mutex
	queues    map[string][]model.CrawlRequest
	inFlight  map[string]bool
	hostOrder []string
}

// NewPolite constructs a per-host scheduler that enforces one in-flight request per host.
func NewPolite() *PoliteFrontier {
	return &PoliteFrontier{
		queues:   make(map[string][]model.CrawlRequest),
		inFlight: make(map[string]bool),
	}
}

// Push enqueues a request into its host-specific FIFO queue.
func (p *PoliteFrontier) Push(req model.CrawlRequest) {
	u, err := url.Parse(req.URL)
	if err != nil || u.Host == "" {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()

	if _, exists := p.queues[u.Host]; !exists {
		p.hostOrder = append(p.hostOrder, u.Host)
	}

	p.queues[u.Host] = append(p.queues[u.Host], req)
}

// Pop selects the next available host (round-robin order) that is not in-flight.
func (p *PoliteFrontier) Pop() (model.CrawlRequest, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, host := range p.hostOrder {
		if p.inFlight[host] {
			continue
		}

		queue := p.queues[host]
		if len(queue) == 0 {
			continue
		}

		req := queue[0]
		p.queues[host] = queue[1:]
		p.inFlight[host] = true
		return req, true
	}

	return model.CrawlRequest{}, false
}

// Done marks a host as no longer in-flight after a request finishes.
func (p *PoliteFrontier) Done(req model.CrawlRequest) {
	u, err := url.Parse(req.URL)
	if err != nil || u.Host == "" {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.inFlight[u.Host] = false
}

// Len returns the total number of queued (not in-flight) requests across all hosts.
func (p *PoliteFrontier) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	total := 0
	for _, q := range p.queues {
		total += len(q)
	}
	return total
}
