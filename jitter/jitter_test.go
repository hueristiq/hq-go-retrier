package jitter_test

import (
	"testing"
	"time"

	"github.com/hueristiq/hq-go-retrier/jitter"
	"github.com/stretchr/testify/assert"
)

func TestEqualJitter(t *testing.T) {
	t.Parallel()

	backoff := 10 * time.Second

	// Run the Equal jitter function multiple times to ensure randomness.
	for range 100 {
		jittered := jitter.Equal(backoff)
		midpoint := backoff / 2

		// Check that jittered value is within the expected range.
		assert.GreaterOrEqual(t, jittered, midpoint, "Jittered duration should be at least the midpoint")
		assert.Less(t, jittered, backoff, "Jittered duration should be less than the original backoff")
	}
}

func TestEqualJitter_MidpointLogic(t *testing.T) {
	t.Parallel()

	backoff := 8 * time.Second
	expectedMidpoint := backoff / 2

	jittered := jitter.Equal(backoff)

	// Check that the jitter is at least the midpoint.
	assert.GreaterOrEqual(t, jittered, expectedMidpoint, "Jittered duration should be at least the midpoint")
}

func TestEqualJitter_ZeroDuration(t *testing.T) {
	t.Parallel()

	backoff := 0 * time.Second

	jittered := jitter.Equal(backoff)

	// Check that when the backoff is 0, the jittered value should also be 0.
	assert.Equal(t, 0*time.Second, jittered, "Jittered duration should be 0 when the backoff is 0")
}

func TestFullJitter(t *testing.T) {
	t.Parallel()

	backoff := 10 * time.Second

	// Run the Full jitter function multiple times to ensure randomness.
	for range 100 {
		jittered := jitter.Full(backoff)

		// Check that jittered value is within the expected range.
		assert.GreaterOrEqual(t, jittered, 0*time.Second, "Jittered duration should be at least 0")
		assert.Less(t, jittered, backoff, "Jittered duration should be less than the original backoff")
	}
}

func TestFullJitter_ZeroDuration(t *testing.T) {
	t.Parallel()

	backoff := 0 * time.Second

	jittered := jitter.Full(backoff)

	// Check that when the backoff is 0, the jittered value should also be 0.
	assert.Equal(t, 0*time.Second, jittered, "Jittered duration should be 0 when the backoff is 0")
}

func TestDecorrelatedJitter_FirstCall(t *testing.T) {
	t.Parallel()

	minDelay := 2 * time.Second
	maxDelay := 10 * time.Second
	previous := 0 * time.Second

	jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

	// Check that jittered value is within the range [minDelay, maxDelay].
	assert.GreaterOrEqual(t, jittered, minDelay, "Jittered duration should be at least the minimum")
	assert.LessOrEqual(t, jittered, maxDelay, "Jittered duration should not exceed the maximum")
}

func TestDecorrelatedJitter_SubsequentCalls(t *testing.T) {
	t.Parallel()

	minDelay := 2 * time.Second
	maxDelay := 10 * time.Second
	previous := 4 * time.Second

	// Run the Decorrelated jitter function multiple times to ensure randomness.
	for range 100 {
		jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

		// Ensure jittered value is within the expected range.
		assert.GreaterOrEqual(t, jittered, minDelay, "Jittered duration should be at least the minimum")
		assert.LessOrEqual(t, jittered, maxDelay, "Jittered duration should not exceed the maximum")
		assert.LessOrEqual(t, jittered, previous*3, "Jittered duration should not exceed three times the previous duration")
	}
}

func TestDecorrelatedJitter_MaxBoundary(t *testing.T) {
	t.Parallel()

	minDelay := 1 * time.Second
	maxDelay := 10 * time.Second
	previous := 5 * time.Second

	jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

	// Ensure that the jittered value is within the expected range and does not exceed the maximum.
	assert.GreaterOrEqual(t, jittered, minDelay, "Jittered duration should be at least the minimum")
	assert.LessOrEqual(t, jittered, maxDelay, "Jittered duration should not exceed the maximum")
}
