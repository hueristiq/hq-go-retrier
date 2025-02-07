package backoff_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.source.hueristiq.com/retrier/backoff"
)

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()

	b := backoff.Exponential()

	tests := []struct {
		minDelay, maxDelay time.Duration
		attempt            int
		expected           time.Duration
	}{
		{time.Millisecond, time.Second, 0, time.Millisecond},                                // i.e 2^0 = 1 * minDelay
		{time.Millisecond, time.Second, 1, 2 * time.Millisecond},                            // i.e 2^1 = 2 * minDelay
		{time.Millisecond, time.Second, 2, 4 * time.Millisecond},                            // i.e 2^2 = 4 * minDelay
		{time.Millisecond, time.Second, 10, time.Second},                                    // Cap at maxDelay
		{50 * time.Millisecond, 2 * time.Second, 5, time.Duration(1600) * time.Millisecond}, // i.e 50 * 2^5
	}

	for _, tt := range tests {
		delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

		assert.Equal(t, tt.expected, delay, "Unexpected backoff duration for attempt %d", tt.attempt)
	}
}

func TestExponentialWithEqualJitterBackoff(t *testing.T) {
	t.Parallel()

	b := backoff.ExponentialWithEqualJitter()

	tests := []struct {
		minDelay, maxDelay time.Duration
		attempt            int
	}{
		{time.Millisecond, time.Second, 0},  // Minimum delay with jitter
		{time.Millisecond, time.Second, 1},  // Jitter added to exponential backoff
		{time.Millisecond, time.Second, 5},  // Check mid-range attempt
		{time.Millisecond, time.Second, 10}, // Maximum delay with jitter
	}

	for _, tt := range tests {
		delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

		assert.GreaterOrEqual(t, delay, tt.minDelay, "Backoff delay should not be less than the minimum")
		assert.LessOrEqual(t, delay, tt.maxDelay, "Backoff delay should not exceed the maximum")
	}
}

func TestExponentialWithFullJitterBackoff(t *testing.T) {
	t.Parallel()

	b := backoff.ExponentialWithFullJitter()

	tests := []struct {
		minDelay, maxDelay time.Duration
		attempt            int
	}{
		{time.Millisecond, time.Second, 0},  // Minimum delay with jitter
		{time.Millisecond, time.Second, 1},  // Jitter added to exponential backoff
		{time.Millisecond, time.Second, 5},  // Check mid-range attempt
		{time.Millisecond, time.Second, 10}, // Maximum delay with jitter
	}

	for _, tt := range tests {
		delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

		assert.GreaterOrEqual(t, delay, tt.minDelay, "Backoff delay should not be less than the minimum")
		assert.LessOrEqual(t, delay, tt.maxDelay, "Backoff delay should not exceed the maximum")
	}
}

func TestExponentialWithDecorrelatedJitterBackoff(t *testing.T) {
	t.Parallel()

	b := backoff.ExponentialWithDecorrelatedJitter()

	tests := []struct {
		minDelay, maxDelay time.Duration
		attempt            int
	}{
		{time.Millisecond, time.Second, 1},  // Decorrelated jitter
		{time.Millisecond, time.Second, 2},  // Jitter added to exponential backoff
		{time.Millisecond, time.Second, 5},  // Check mid-range attempt
		{time.Millisecond, time.Second, 10}, // Maximum delay with decorrelated jitter
	}

	for _, tt := range tests {
		delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

		assert.GreaterOrEqual(t, delay, tt.minDelay, "Backoff delay should not be less than the minimum")
		assert.LessOrEqual(t, delay, tt.maxDelay, "Backoff delay should not exceed the maximum")
	}
}
