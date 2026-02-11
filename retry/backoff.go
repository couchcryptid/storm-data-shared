package retry

import (
	"context"
	"time"
)

// NextBackoff doubles the current backoff, capping at maxBackoff.
func NextBackoff(current, maxBackoff time.Duration) time.Duration {
	next := current * 2
	if next > maxBackoff {
		return maxBackoff
	}
	return next
}

// SleepWithContext sleeps for d, returning false if the context is cancelled first.
// Returns true immediately for zero or negative durations.
func SleepWithContext(ctx context.Context, d time.Duration) bool {
	if d <= 0 {
		return true
	}

	timer := time.NewTimer(d)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
