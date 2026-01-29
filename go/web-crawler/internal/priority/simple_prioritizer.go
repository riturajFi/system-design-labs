package priority

import (
	"strings"
	"web-crawler/internal/model"
)

// SimplePrioritizer uses a basic substring heuristic to assign priority.
type SimplePrioritizer struct{}

// NewSimple creates a replaceable, minimal prioritizer implementation.
func NewSimple() *SimplePrioritizer {
	return &SimplePrioritizer{}
}

// Assign sets higher priority for URLs matching a known substring.
func (p *SimplePrioritizer) Assign(req model.CrawlRequest) model.CrawlRequest {
	if strings.Contains(req.URL, "example.com") {
		req.Priority = model.PriorityHigh
		return req
	}
	req.Priority = model.PriorityLow
	return req
}
