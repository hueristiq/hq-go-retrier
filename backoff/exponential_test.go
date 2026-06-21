package backoff

import (
	"math"
	"testing"
	"time"

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

		b := Exponential()

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

		b := Exponential()

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
		b := Exponential()

		delay := b(minDelay, maxDelay, 2)

		assert.Equal(t, maxDelay, delay, "Should cap at maxDelay when overflow would occur")
	})
}

func expBase(minDelay, maxDelay time.Duration, attempt int) time.Duration {
	return min(minDelay<<attempt, maxDelay)
}

func TestExponentialDeterminism(t *testing.T) {
	t.Parallel()

	b := Exponential()
	want := b(time.Millisecond, time.Second, 4)

	for range 20 {
		assert.Equal(t, want, b(time.Millisecond, time.Second, 4), "Exponential should be deterministic")
	}
}

func TestExponentialJitterVaries(t *testing.T) {
	t.Parallel()

	const samples = 50

	strategies := []struct {
		name string
		b    Backoff
	}{
		{"equal jitter", ExponentialWithEqualJitter()},
		{"full jitter", ExponentialWithFullJitter()},
		{"decorrelated jitter", ExponentialWithDecorrelatedJitter()},
	}

	for _, s := range strategies {
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			seen := make(map[time.Duration]struct{})

			for range samples {
				seen[s.b(time.Second, time.Minute, 3)] = struct{}{}
			}

			assert.Greater(t, len(seen), 1, "Jittered strategy should produce varied delays")
		})
	}
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
			{"attempt 1", time.Millisecond, time.Second, 1},
			{"attempt 2", time.Millisecond, time.Second, 2},
			{"attempt 3", time.Millisecond, time.Second, 3},
			{"attempt 4", time.Millisecond, time.Second, 4},
			{"attempt 5", time.Millisecond, time.Second, 5},
			{"attempt 6", time.Millisecond, time.Second, 6},
			{"attempt 7", time.Millisecond, time.Second, 7},
			{"attempt 8", time.Millisecond, time.Second, 8},
			{"attempt 9", time.Millisecond, time.Second, 9},
			{"attempt 10 (capped)", time.Millisecond, 2 * time.Second, 10},
		}

		b := ExponentialWithEqualJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				base := expBase(tt.minDelay, tt.maxDelay, tt.attempt)

				for range 10 {
					delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

					assert.GreaterOrEqual(t, delay, base/2, "Delay should be at least base/2")
					assert.LessOrEqual(t, delay, base, "Delay should not exceed base")
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
			expectZero         bool
		}{
			{"negative minDelay", -time.Millisecond, time.Second, 1, true},
			{"negative maxDelay", time.Millisecond, -time.Second, 1, true},
			{"negative attempt", time.Millisecond, time.Second, -1, true},
			{"minDelay = maxDelay", time.Second, time.Second, 5, false},
			{"minDelay > maxDelay", 2 * time.Second, time.Second, 0, false},
			{"zero attempt", time.Millisecond, time.Second, 0, false},
		}

		b := ExponentialWithEqualJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

				if tt.expectZero {
					assert.Equal(t, time.Duration(0), delay)

					return
				}

				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, tt.maxDelay)
			})
		}
	})

	t.Run("overflow protection", func(t *testing.T) {
		t.Parallel()

		minDelay := time.Duration(math.MaxInt64 / 2)
		maxDelay := time.Duration(math.MaxInt64)
		b := ExponentialWithEqualJitter()

		delay := b(minDelay, maxDelay, 2)

		assert.GreaterOrEqual(t, delay, time.Duration(0))
		assert.LessOrEqual(t, delay, maxDelay, "Should not exceed maxDelay when overflow would occur")
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
			{"attempt 1", time.Millisecond, time.Second, 1},
			{"attempt 2", time.Millisecond, time.Second, 2},
			{"attempt 3", time.Millisecond, time.Second, 3},
			{"attempt 4", time.Millisecond, time.Second, 4},
			{"attempt 5", time.Millisecond, time.Second, 5},
			{"attempt 6", time.Millisecond, time.Second, 6},
			{"attempt 7", time.Millisecond, time.Second, 7},
			{"attempt 8", time.Millisecond, time.Second, 8},
			{"attempt 9", time.Millisecond, time.Second, 9},
			{"attempt 10 (capped)", time.Millisecond, time.Second, 10},
		}

		b := ExponentialWithFullJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				base := expBase(tt.minDelay, tt.maxDelay, tt.attempt)

				for range 10 {
					delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

					assert.GreaterOrEqual(t, delay, time.Duration(0), "Delay should be at least 0")
					assert.LessOrEqual(t, delay, base, "Delay should not exceed base")
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
			expectZero         bool
		}{
			{"negative minDelay", -time.Millisecond, time.Second, 1, true},
			{"negative maxDelay", time.Millisecond, -time.Second, 1, true},
			{"negative attempt", time.Millisecond, time.Second, -1, true},
			{"minDelay = maxDelay", time.Second, time.Second, 5, false},
			{"minDelay > maxDelay", 2 * time.Second, time.Second, 0, false},
			{"zero attempt", time.Millisecond, time.Second, 0, false},
		}

		b := ExponentialWithFullJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

				if tt.expectZero {
					assert.Equal(t, time.Duration(0), delay)

					return
				}

				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, tt.maxDelay)
			})
		}
	})

	t.Run("overflow protection", func(t *testing.T) {
		t.Parallel()

		minDelay := time.Duration(math.MaxInt64 / 2)
		maxDelay := time.Duration(math.MaxInt64)
		b := ExponentialWithFullJitter()

		delay := b(minDelay, maxDelay, 2)

		assert.GreaterOrEqual(t, delay, time.Duration(0))
		assert.LessOrEqual(t, delay, maxDelay, "Should not exceed maxDelay when overflow would occur")
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
			{"attempt 1", time.Millisecond, time.Second, 1},
			{"attempt 2", time.Millisecond, time.Second, 2},
			{"attempt 3", time.Millisecond, time.Second, 3},
			{"attempt 4", time.Millisecond, time.Second, 4},
			{"attempt 5", time.Millisecond, time.Second, 5},
			{"attempt 6", time.Millisecond, time.Second, 6},
			{"attempt 7", time.Millisecond, time.Second, 7},
			{"attempt 8", time.Millisecond, time.Second, 8},
			{"attempt 9", time.Millisecond, time.Second, 9},
			{"attempt 10 (capped)", time.Millisecond, time.Second, 10},
		}

		b := ExponentialWithDecorrelatedJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				for range 10 {
					delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

					assert.GreaterOrEqual(t, delay, tt.minDelay, "Delay should be at least minDelay")
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
			expectZero         bool
		}{
			{"negative minDelay", -time.Millisecond, time.Second, 1, true},
			{"negative maxDelay", time.Millisecond, -time.Second, 1, true},
			{"negative attempt", time.Millisecond, time.Second, -1, true},
			{"minDelay > maxDelay", 2 * time.Second, time.Second, 0, true},
			{"minDelay = maxDelay", time.Second, time.Second, 5, false},
			{"zero attempt", time.Millisecond, time.Second, 0, false},
		}

		b := ExponentialWithDecorrelatedJitter()

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				t.Parallel()

				delay := b(tt.minDelay, tt.maxDelay, tt.attempt)

				if tt.expectZero {
					assert.Equal(t, time.Duration(0), delay)

					return
				}

				assert.GreaterOrEqual(t, delay, time.Duration(0))
				assert.LessOrEqual(t, delay, tt.maxDelay)
			})
		}
	})

	t.Run("overflow protection", func(t *testing.T) {
		t.Parallel()

		minDelay := time.Duration(math.MaxInt64 / 2)
		maxDelay := time.Duration(math.MaxInt64)
		b := ExponentialWithDecorrelatedJitter()

		delay := b(minDelay, maxDelay, 2)

		assert.GreaterOrEqual(t, delay, time.Duration(0))
		assert.LessOrEqual(t, delay, maxDelay, "Should not exceed maxDelay when overflow would occur")
	})
}
