package fixed

import (
	"time"

	"notification-system/internal/core/ports"
)

type Policy struct {
	maxAttempts int
	baseDelay   time.Duration
}

func New(maxAttempts int, baseDelay time.Duration) *Policy {
	return &Policy{
		maxAttempts: maxAttempts,
		baseDelay:   baseDelay,
	}
}

func (p *Policy) Decide(
	attempt int,
	_ error,
) ports.RetryDecision {
	if attempt >= p.maxAttempts {
		return ports.RetryDecision{
			Retry: false,
		}
	}

	return ports.RetryDecision{
		Retry: true,
		After: time.Duration(attempt+1) * p.baseDelay,
	}
}
