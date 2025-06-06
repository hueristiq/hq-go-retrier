package jitter_test

import (
	"testing"
	"time"

	"github.com/hueristiq/hq-go-retrier/jitter"
	"github.com/stretchr/testify/assert"
)

func TestEqualJitter(t *testing.T) {
	t.Parallel()

	t.Run("negative backoff", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, time.Duration(0), jitter.Equal(-time.Second))
	})

	t.Run("zero backoff", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, time.Duration(0), jitter.Equal(0))
	})

	t.Run("positive backoff", func(t *testing.T) {
		t.Parallel()

		backoff := 10 * time.Second

		for range 100 {
			jittered := jitter.Equal(backoff)
			midpoint := backoff / 2

			assert.GreaterOrEqual(t, jittered, midpoint, "Jittered duration should be at least the midpoint")
			assert.LessOrEqual(t, jittered, backoff, "Jittered duration should not exceed the original backoff")
		}
	})

	t.Run("small backoff", func(t *testing.T) {
		t.Parallel()

		backoff := 1 * time.Nanosecond
		jittered := jitter.Equal(backoff)

		assert.Equal(t, backoff/2, jittered, "For very small backoffs, should return midpoint")
	})
}

func TestFullJitter(t *testing.T) {
	t.Parallel()

	t.Run("negative backoff", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, time.Duration(0), jitter.Full(-time.Second))
	})

	t.Run("zero backoff", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, time.Duration(0), jitter.Full(0))
	})

	t.Run("positive backoff", func(t *testing.T) {
		t.Parallel()

		backoff := 10 * time.Second

		for range 100 {
			jittered := jitter.Full(backoff)

			assert.GreaterOrEqual(t, jittered, 0*time.Second, "Jittered duration should be at least 0")
			assert.Less(t, jittered, backoff, "Jittered duration should be less than the original backoff")
		}
	})

	t.Run("small backoff", func(t *testing.T) {
		t.Parallel()

		backoff := 1 * time.Nanosecond
		jittered := jitter.Full(backoff)

		assert.Equal(t, time.Duration(0), jittered, "For very small backoffs, should return 0")
	})
}
func TestDecorrelatedJitter(t *testing.T) {
	t.Parallel()

	t.Run("negative min/max", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, time.Duration(0), jitter.Decorrelated(-1, 10, 0))
		assert.Equal(t, time.Duration(0), jitter.Decorrelated(1, -10, 0))
		assert.Equal(t, time.Duration(0), jitter.Decorrelated(-1, -10, 0))
	})

	t.Run("min > max", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, time.Duration(0), jitter.Decorrelated(10, 5, 0))
	})

	t.Run("min equals max", func(t *testing.T) {
		t.Parallel()

		delay := 5 * time.Second
		jittered := jitter.Decorrelated(delay, delay, 0)
		assert.Equal(t, delay, jittered)
	})

	t.Run("first call (previous=0)", func(t *testing.T) {
		t.Parallel()

		minDelay := 2 * time.Second
		maxDelay := 10 * time.Second
		previous := 0 * time.Second

		jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

		assert.GreaterOrEqual(t, jittered, minDelay, "Jittered duration should be at least the minimum")
		assert.LessOrEqual(t, jittered, maxDelay, "Jittered duration should not exceed the maximum")
	})

	t.Run("subsequent calls", func(t *testing.T) {
		t.Parallel()

		minDelay := 2 * time.Second
		maxDelay := 10 * time.Second
		previous := 4 * time.Second

		for range 100 {
			jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

			assert.GreaterOrEqual(t, jittered, minDelay, "Jittered duration should be at least the minimum")
			assert.LessOrEqual(t, jittered, maxDelay, "Jittered duration should not exceed the maximum")
			assert.LessOrEqual(t, jittered, previous*3+minDelay, "Jittered duration should not exceed three times the previous duration plus minDelay")
		}
	})

	t.Run("max boundary", func(t *testing.T) {
		t.Parallel()

		minDelay := 1 * time.Second
		maxDelay := 10 * time.Second
		previous := 5 * time.Second

		// Run multiple times to catch potential boundary issues
		for range 100 {
			jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

			assert.GreaterOrEqual(t, jittered, minDelay, "Jittered duration should be at least the minimum")
			assert.LessOrEqual(t, jittered, maxDelay, "Jittered duration should not exceed the maximum")
		}
	})

	t.Run("previous causes overflow", func(t *testing.T) {
		t.Parallel()

		minDelay := 1 * time.Second
		maxDelay := 10 * time.Second
		previous := time.Duration(1<<63 - 1) // Max possible duration

		jittered := jitter.Decorrelated(minDelay, maxDelay, previous)

		assert.Equal(t, maxDelay, jittered)
	})
}
