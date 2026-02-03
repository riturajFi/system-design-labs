package metrics

import (
	"sync"
)

type Registry struct {
	mu             sync.Mutex
	healthChecks   uint64
	rateAllowed    uint64
	rateDenied     uint64
	queueDepth     map[string]int
	workerDequeued uint64
	trackingEmitted uint64
	trackingFailed  uint64
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

func (r *Registry) SetQueueDepth(name string, depth int) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.queueDepth == nil {
		r.queueDepth = make(map[string]int)
	}
	r.queueDepth[name] = depth
}

func (r *Registry) IncWorkerDequeued() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.workerDequeued++
}

func (r *Registry) IncTrackingEmitted() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trackingEmitted++
}

func (r *Registry) IncTrackingFailed() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.trackingFailed++
}

func (r *Registry) Snapshot() map[string]uint64 {
	r.mu.Lock()
	defer r.mu.Unlock()

	snap := map[string]uint64{
		"health_checks_total":    r.healthChecks,
		"rate_allowed_total":     r.rateAllowed,
		"rate_denied_total":      r.rateDenied,
		"worker_dequeued_total":  r.workerDequeued,
		"tracking_emitted_total": r.trackingEmitted,
		"tracking_failed_total":  r.trackingFailed,
	}

	for k, v := range r.queueDepth {
		snap["queue_depth_"+k] = uint64(v)
	}

	return snap
}
