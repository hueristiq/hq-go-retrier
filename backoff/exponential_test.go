package backoff_test

import (
	"math"
	"testing"
	"time"

	"github.com/hueristiq/hq-go-retrier/backoff"
	"github.com/stretchr/testify/assert"
)

func TestExponentialBackoff(t *testing.T) {
	t.Parallel()

	t.Run("standard progression", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name               string
			minDelay, maxDelay time.Duration
			attempt            int
			expected           time.Duration
		}{
			{
				name:     "attempt 1",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  1,
				expected: 2 * time.Millisecond,
			},
			{
				name:     "attempt 2",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  2,
				expected: 4 * time.Millisecond,
			},
			{
				name:     "attempt 3",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  3,
				expected: 8 * time.Millisecond,
			},
			{
				name:     "attempt 4",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  4,
				expected: 16 * time.Millisecond,
			},
			{
				name:     "attempt 5",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  5,
				expected: 32 * time.Millisecond,
			},
			{
				name:     "attempt 6",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  6,
				expected: 64 * time.Millisecond,
			},
			{
				name:     "attempt 7",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  7,
				expected: 128 * time.Millisecond,
			},
			{
				name:     "attempt 8",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  8,
				expected: 256 * time.Millisecond,
			},
			{
				name:     "attempt 9",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  9,
				expected: 512 * time.Millisecond,
			},
			{
				name:     "attempt 10 (capped)",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  10,
				expected: time.Second,
			},
		}

		b := backoff.Exponential()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

				assert.Equal(t, tt.expected, delay, "Unexpected backoff duration for attempt %d", tt.attempt)
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name               string
			minDelay, maxDelay time.Duration
			attempt            int
			expected           time.Duration
		}{
			{
				name:     "negative minDelay",
				minDelay: -time.Millisecond,
				maxDelay: time.Second,
				attempt:  1,
				expected: 0,
			},
			{
				name:     "negative maxDelay",
				minDelay: time.Millisecond,
				maxDelay: -time.Second,
				attempt:  1,
				expected: 0,
			},
			{
				name:     "minDelay = maxDelay",
				minDelay: time.Second,
				maxDelay: time.Second,
				attempt:  5,
				expected: time.Second,
			},
			{
				name:     "minDelay > maxDelay",
				minDelay: 2 * time.Second,
				maxDelay: time.Second,
				attempt:  0,
				expected: time.Second,
			},
			{
				name:     "negative attempt",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  -1,
				expected: 0,
			},
			{
				name:     "zero attempt",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  0,
				expected: time.Millisecond,
			},
		}

		b := backoff.Exponential()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

				assert.Equal(t, tt.expected, delay)
			})
		}
	})

	t.Run("overflow protection", func(t *testing.T) {
		t.Parallel()

		minDelay := time.Duration(math.MaxInt64 / 2)
		maxDelay := time.Duration(math.MaxInt64)
		b := backoff.Exponential()

		delay := b(minDelay, maxDelay, 2)

		assert.Equal(t, maxDelay, delay, "Should cap at maxDelay when overflow would occur")
	})
}

func TestExponentialWithEqualJitterBackoff(t *testing.T) {
	t.Parallel()

	t.Run("jitter range validation", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name               string
			minDelay, maxDelay time.Duration
			attempt            int
		}{
			{
				name:     "attempt 1",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  1,
			},
			{
				name:     "attempt 2",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  2,
			},
			{
				name:     "attempt 3",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  3,
			},
			{
				name:     "attempt 4",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  4,
			},
			{
				name:     "attempt 5",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  5,
			},
			{
				name:     "attempt 6",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  6,
			},
			{
				name:     "attempt 7",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  7,
			},
			{
				name:     "attempt 8",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  8,
			},
			{
				name:     "attempt 9",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  9,
			},
			{
				name:     "attempt 10 (capped)",
				minDelay: time.Millisecond,
				maxDelay: 2 * time.Second,
				attempt:  10,
			},
		}

		b := backoff.ExponentialWithEqualJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				base := tt.minDelay << tt.attempt
				if base > tt.maxDelay {
					base = tt.maxDelay
				}

				midpoint := base / 2

				for range 10 {
					delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

					assert.GreaterOrEqual(t, delay, base+midpoint, "Delay should be at least base + midpoint")
					assert.LessOrEqual(t, delay, base*2, "Delay should not exceed base * 2 backoff")
					assert.LessOrEqual(t, delay, tt.maxDelay, "Delay should not exceed maxDelay")
				}
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name               string
			minDelay, maxDelay time.Duration
			attempt            int
			expected           time.Duration
		}{
			{
				name:     "negative minDelay",
				minDelay: -time.Millisecond,
				maxDelay: time.Second,
				attempt:  1,
				expected: 0,
			},
			{
				name:     "negative maxDelay",
				minDelay: time.Millisecond,
				maxDelay: -time.Second,
				attempt:  1,
				expected: 0,
			},
			{
				name:     "minDelay = maxDelay",
				minDelay: time.Second,
				maxDelay: time.Second,
				attempt:  5,
				expected: time.Second,
			},
			{
				name:     "minDelay > maxDelay",
				minDelay: 2 * time.Second,
				maxDelay: time.Second,
				attempt:  0,
				expected: time.Second,
			},
			{
				name:     "negative attempt",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  -1,
				expected: 0,
			},
			{
				name:     "zero attempt",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  0,
				expected: time.Millisecond,
			},
		}

		b := backoff.ExponentialWithEqualJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

				assert.Equal(t, tt.expected, delay)
			})
		}
	})

	t.Run("overflow protection", func(t *testing.T) {
		t.Parallel()

		minDelay := time.Duration(math.MaxInt64 / 2)
		maxDelay := time.Duration(math.MaxInt64)
		b := backoff.ExponentialWithEqualJitter()

		delay := b(minDelay, maxDelay, 2)

		assert.Equal(t, maxDelay, delay, "Should cap at maxDelay when overflow would occur")
	})
}

func TestExponentialWithFullJitterBackoff(t *testing.T) {
	t.Parallel()

	t.Run("jitter range validation", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name               string
			minDelay, maxDelay time.Duration
			attempt            int
		}{
			{
				name:     "attempt 1",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  1,
			},
			{
				name:     "attempt 2",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  2,
			},
			{
				name:     "attempt 3",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  3,
			},
			{
				name:     "attempt 4",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  4,
			},
			{
				name:     "attempt 5",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  5,
			},
			{
				name:     "attempt 6",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  6,
			},
			{
				name:     "attempt 7",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  7,
			},
			{
				name:     "attempt 8",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  8,
			},
			{
				name:     "attempt 9",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  9,
			},
			{
				name:     "attempt 10 (capped)",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  10,
			},
		}

		b := backoff.ExponentialWithFullJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				base := tt.minDelay << tt.attempt
				if base > tt.maxDelay {
					base = tt.maxDelay
				}

				for range 10 {
					delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

					assert.GreaterOrEqual(t, delay, base, "Delay should be at least base")
					assert.LessOrEqual(t, delay, base*2, "Delay should not exceed base * 2 backoff")
					assert.LessOrEqual(t, delay, tt.maxDelay, "Delay should not exceed maxDelay")
				}
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name               string
			minDelay, maxDelay time.Duration
			attempt            int
			expected           time.Duration
		}{

			{
				name:     "negative minDelay",
				minDelay: -time.Millisecond,
				maxDelay: time.Second,
				attempt:  1,
				expected: 0,
			},
			{
				name:     "negative maxDelay",
				minDelay: time.Millisecond,
				maxDelay: -time.Second,
				attempt:  1,
				expected: 0,
			},
			{
				name:     "minDelay = maxDelay",
				minDelay: time.Second,
				maxDelay: time.Second,
				attempt:  5,
				expected: time.Second,
			},
			{
				name:     "minDelay > maxDelay",
				minDelay: 2 * time.Second,
				maxDelay: time.Second,
				attempt:  0,
				expected: time.Second,
			},
			{
				name:     "negative attempt",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  -1,
				expected: 0,
			},
			{
				name:     "zero attempt",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  0,
				expected: time.Millisecond,
			},
		}

		b := backoff.ExponentialWithFullJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

				assert.Equal(t, tt.expected, delay)
			})
		}
	})

	t.Run("overflow protection", func(t *testing.T) {
		t.Parallel()

		minDelay := time.Duration(math.MaxInt64 / 2)
		maxDelay := time.Duration(math.MaxInt64)
		b := backoff.ExponentialWithFullJitter()

		delay := b(minDelay, maxDelay, 2)

		assert.Equal(t, maxDelay, delay, "Should cap at maxDelay when overflow would occur")
	})
}

func TestExponentialWithDecorrelatedJitterBackoff(t *testing.T) {
	t.Parallel()

	t.Run("jitter progression", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name               string
			minDelay, maxDelay time.Duration
			attempt            int
		}{
			{
				name:     "attempt 1",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  1,
			},
			{
				name:     "attempt 2",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  2,
			},
			{
				name:     "attempt 3",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  3,
			},
			{
				name:     "attempt 4",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  4,
			},
			{
				name:     "attempt 5",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  5,
			},
			{
				name:     "attempt 6",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  6,
			},
			{
				name:     "attempt 7",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  7,
			},
			{
				name:     "attempt 8",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  8,
			},
			{
				name:     "attempt 9",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  9,
			},
			{
				name:     "attempt 10 (capped)",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  10,
			},
		}

		b := backoff.ExponentialWithDecorrelatedJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				base := tt.minDelay << tt.attempt
				if base > tt.maxDelay {
					base = tt.maxDelay
				}

				previous := tt.minDelay

				if tt.attempt > 0 {
					previous = tt.minDelay << (tt.attempt - 1)
				}

				for range 10 {
					delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

					assert.GreaterOrEqual(t, delay, base, "Delay should be at least base")
					assert.LessOrEqual(t, delay, base+(tt.minDelay+(previous*3)), "Delay should not exceed base+(tt.minDelay+(previous*3)) backoff")
					assert.LessOrEqual(t, delay, tt.maxDelay, "Delay should not exceed maxDelay")
				}
			})
		}
	})

	t.Run("edge cases", func(t *testing.T) {
		t.Parallel()

		tests := []struct {
			name               string
			minDelay, maxDelay time.Duration
			attempt            int
			expected           time.Duration
		}{
			{
				name:     "negative minDelay",
				minDelay: -time.Millisecond,
				maxDelay: time.Second,
				attempt:  1,
				expected: 0,
			},
			{
				name:     "negative maxDelay",
				minDelay: time.Millisecond,
				maxDelay: -time.Second,
				attempt:  1,
				expected: 0,
			},
			{
				name:     "minDelay = maxDelay",
				minDelay: time.Second,
				maxDelay: time.Second,
				attempt:  5,
				expected: time.Second,
			},
			{
				name:     "minDelay > maxDelay",
				minDelay: 2 * time.Second,
				maxDelay: time.Second,
				attempt:  0,
				expected: time.Second,
			},
			{
				name:     "negative attempt",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  -1,
				expected: 0,
			},
			{
				name:     "zero attempt",
				minDelay: time.Millisecond,
				maxDelay: time.Second,
				attempt:  0,
				expected: time.Millisecond,
			},
		}

		b := backoff.ExponentialWithDecorrelatedJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

				assert.Equal(t, tt.expected, delay)
			})
		}
	})

	t.Run("overflow protection", func(t *testing.T) {
		t.Parallel()

		minDelay := time.Duration(math.MaxInt64 / 2)
		maxDelay := time.Duration(math.MaxInt64)
		b := backoff.ExponentialWithDecorrelatedJitter()

		delay := b(minDelay, maxDelay, 2)

		assert.Equal(t, maxDelay, delay, "Should cap at maxDelay when overflow would occur")
	})
}
