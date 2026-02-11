package retry

import (
	"context"
	"testing"
	"time"
)

func TestNextBackoff(t *testing.T) {
	tests := []struct {
		name    string
		current time.Duration
		max     time.Duration
		want    time.Duration
	}{
		{"doubles", 200 * time.Millisecond, 5 * time.Second, 400 * time.Millisecond},
		{"doubles again", 400 * time.Millisecond, 5 * time.Second, 800 * time.Millisecond},
		{"caps at max", 4 * time.Second, 5 * time.Second, 5 * time.Second},
		{"stays at max", 5 * time.Second, 5 * time.Second, 5 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextBackoff(tt.current, tt.max)
			if got != tt.want {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSleepWithContext_Completes(t *testing.T) {
	ctx := context.Background()
	if !SleepWithContext(ctx, 1*time.Millisecond) {
		t.Error("expected true for completed sleep")
	}
}

func TestSleepWithContext_ContextCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if SleepWithContext(ctx, 1*time.Second) {
		t.Error("expected false for cancelled context")
	}
}

func TestSleepWithContext_ZeroDuration(t *testing.T) {
	ctx := context.Background()
	if !SleepWithContext(ctx, 0) {
		t.Error("expected true for zero duration")
	}
}
