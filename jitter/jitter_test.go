package jitter_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.sources.hueristiq.com/retrier/jitter"
)

func TestEqualJitter(t *testing.T) {
	t.Parallel()

	backoff := 10 * time.Second

	for range 100 {
		jittered := jitter.Equal(backoff)
		midpoint := backoff / 2

		assert.GreaterOrEqual(t, jittered, midpoint, "Jittered duration should be at least the midpoint")
		assert.Less(t, jittered, backoff, "Jittered duration should be less than the original backoff")
	}
}

func TestEqualJitter_MidpointLogic(t *testing.T) {
	t.Parallel()

	backoff := 8 * time.Second
	expectedMidpoint := backoff / 2

	jittered := jitter.Equal(backoff)

	assert.GreaterOrEqual(t, jittered, expectedMidpoint, "Jittered duration should be at least the midpoint")
}

func TestEqualJitter_ZeroDuration(t *testing.T) {
	t.Parallel()

	backoff := 0 * time.Second

	jittered := jitter.Equal(backoff)

	assert.Equal(t, 0*time.Second, jittered, "Jittered duration should be 0 when the backoff is 0")
}

func TestFullJitter(t *testing.T) {
	t.Parallel()

	backoff := 10 * time.Second

	for range 100 {
		jittered := jitter.Full(backoff)

		assert.GreaterOrEqual(t, jittered, 0*time.Second, "Jittered duration should be at least 0")
		assert.Less(t, jittered, backoff, "Jittered duration should be less than the original backoff")
	}
}

func TestFullJitter_ZeroDuration(t *testing.T) {
	t.Parallel()

	backoff := 0 * time.Second

	jittered := jitter.Full(backoff)

	assert.Equal(t, 0*time.Second, jittered, "Jittered duration should be 0 when the backoff is 0")
}

func TestDecorrelatedJitter_FirstCall(t *testing.T) {
	t.Parallel()

	minDelay := 2 * time.Second
	maxDelay := 10 * time.Second
	previous := 0 * time.Second

	jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

	assert.GreaterOrEqual(t, jittered, minDelay, "Jittered duration should be at least the minimum")
	assert.LessOrEqual(t, jittered, maxDelay, "Jittered duration should not exceed the maximum")
}

func TestDecorrelatedJitter_SubsequentCalls(t *testing.T) {
	t.Parallel()

	minDelay := 2 * time.Second
	maxDelay := 10 * time.Second
	previous := 4 * time.Second

	for range 100 {
		jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

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

	assert.GreaterOrEqual(t, jittered, minDelay, "Jittered duration should be at least the minimum")
	assert.LessOrEqual(t, jittered, maxDelay, "Jittered duration should not exceed the maximum")
}
