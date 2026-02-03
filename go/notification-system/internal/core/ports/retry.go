package ports

import "time"

type RetryDecision struct {
	Retry bool
	After time.Duration
}

type RetryPolicy interface {
	Decide(attempt int, err error) RetryDecision
}
