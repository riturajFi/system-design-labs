package metrics

import (
	"sync"
)

type Registry struct {
	mu           sync.Mutex
	healthChecks uint64
	rateAllowed  uint64
	rateDenied   uint64
}

func NewRegistry() *Registry {
	return &Registry{}
}

func (r *Registry) IncHealthCheck() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.healthChecks++
}

func (r *Registry) IncRateAllowed() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rateAllowed++
}

func (r *Registry) IncRateDenied() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.rateDenied++
}

func (r *Registry) Snapshot() map[string]uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	return map[string]uint64{
		"health_checks_total": r.healthChecks,
		"rate_allowed_total":  r.rateAllowed,
		"rate_denied_total":   r.rateDenied,
	}
}
