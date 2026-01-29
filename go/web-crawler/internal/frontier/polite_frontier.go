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

func NewPolite() *PoliteFrontier {
	return &PoliteFrontier{
		queues:   make(map[string][]model.CrawlRequest),
		inFlight: make(map[string]bool),
	}
}

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

func (p *PoliteFrontier) Done(req model.CrawlRequest) {
	u, err := url.Parse(req.URL)
	if err != nil || u.Host == "" {
		return
	}

	p.mu.Lock()
	defer p.mu.Unlock()
	p.inFlight[u.Host] = false
}

func (p *PoliteFrontier) Len() int {
	p.mu.Lock()
	defer p.mu.Unlock()

	total := 0
	for _, q := range p.queues {
		total += len(q)
	}
	return total
}
