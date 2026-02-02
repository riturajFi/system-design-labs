package metrics

import (
	"sync"
)

type Registry struct {
	mu sync.Mutex
	healthChecks uint64
}

func NewRegistry() *Registry {
	return  &Registry{}
}

func (r *Registry) IncHealthCheck() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.healthChecks ++
}

func(r *Registry) Snapshot() map[string]uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	return map[string]uint64{
		"health_checks_total": r.healthChecks,
	}
}