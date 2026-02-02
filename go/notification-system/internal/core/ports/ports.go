package ports

// Metrics is a minimal interface for emitting metrics without coupling core logic
// to a specific metrics backend.
type Metrics interface {
	IncCounter(name string, tags map[string]string)
}

type NoopMetrics struct{}

func (NoopMetrics) IncCounter(_ string, _ map[string]string) {}
