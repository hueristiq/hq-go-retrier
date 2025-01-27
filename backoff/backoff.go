package backoff

import (
	"math"
	"sync"
	"time"

	"go.sources.hueristiq.com/retrier/jitter"
)

// Backoff is a function type that calculates the delay between retry attempts.
// It takes a minimum duration, maximum duration, and the current retry attempt number as inputs,
// and returns the calculated delay duration. Each retry strategy returns a function of this type.
//
// Parameters:
//   - minDelay: The minimum allowable delay duration.
//   - maxDelay: The maximum allowable delay duration.
//   - attempt:  The current retry attempt number.
//
// Returns:
//   - delay: The delay duration to wait before the next retry.
type Backoff func(minDelay, maxDelay time.Duration, attempt int) (delay time.Duration)

// Exponential returns a backoff function that implements basic exponential backoff.
// In this strategy, the delay increases exponentially with each retry attempt,
// but is capped by the provided maximum duration.
//
// Formula: delay = minDelay * 2^attempt
//
// Parameters:
//   - minDelay: The minimum backoff duration (base duration).
//   - maxDelay: The maximum allowable backoff duration.
//   - attempt:  The current retry attempt number.
//
// Returns:
//   - delay: The calculated delay duration, capped at the maximum duration.
//
// Example:
//
//	backoffFunc := backoff.Exponential()
//	delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	// delay will be 8 seconds (1s * 2^3), but capped at maxDelay if exceeded.
func Exponential() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = time.Duration(math.Pow(2, float64(attempt)) * float64(minDelay))

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithEqualJitter returns a backoff function that implements exponential backoff with equal jitter.
// In this strategy, the base delay increases exponentially, and equal jitter is applied to introduce
// moderate randomness by adding a random value from the midpoint of the calculated delay.
//
// Formula: delay = minDelay * 2^attempt + random(midpoint, delay)
//
// Parameters:
//   - minDelay: The minimum backoff duration (base duration).
//   - maxDelay: The maximum allowable backoff duration.
//   - attempt:  The current retry attempt number.
//
// Returns:
//   - delay: The calculated delay with equal jitter applied, capped at the maximum duration.
//
// Example:
//
//	backoffFunc := backoff.ExponentialWithEqualJitter()
//	delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	// delay will be exponentially calculated with equal jitter applied.
func ExponentialWithEqualJitter() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	mutex := &sync.Mutex{}

	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = time.Duration(math.Pow(2, float64(attempt)) * float64(minDelay))

		mutex.Lock()
		jittered := jitter.Equal(backoff)
		mutex.Unlock()

		backoff += jittered

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithFullJitter returns a backoff function that implements exponential backoff with full jitter.
// In this strategy, the base delay increases exponentially, and full jitter is applied by adding a random value
// drawn from a uniform distribution between 0 and the calculated delay.
//
// Formula: delay = minDelay * 2^attempt + random(0, delay)
//
// Parameters:
//   - minDelay: The minimum backoff duration (base duration).
//   - maxDelay: The maximum allowable backoff duration.
//   - attempt:  The current retry attempt number.
//
// Returns:
//   - delay: The calculated delay with full jitter applied, capped at the maximum duration.
//
// Example:
//
//	backoffFunc := backoff.ExponentialWithFullJitter()
//	delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	// delay will be exponentially calculated with full jitter applied.
func ExponentialWithFullJitter() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	mutex := &sync.Mutex{}

	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = time.Duration(math.Pow(2, float64(attempt)) * float64(minDelay))

		mutex.Lock()
		jittered := jitter.Full(backoff)
		mutex.Unlock()

		backoff += jittered

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithDecorrelatedJitter returns a backoff function that implements exponential backoff
// with decorrelated jitter. The base delay increases exponentially, and a decorrelated jitter is applied,
// where the jittered value is influenced by the previous backoff duration.
//
// Formula: delay = minDelay * 2^attempt + random(previous * 3, delay)
//
// Parameters:
//   - minDelay: The minimum backoff duration (base duration).
//   - maxDelay: The maximum allowable backoff duration.
//   - attempt:  The current retry attempt number.
//
// Returns:
//   - delay: The calculated delay with decorrelated jitter applied, capped at the maximum duration.
//
// Example:
//
//	backoffFunc := backoff.ExponentialWithDecorrelatedJitter()
//	delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	// delay will be exponentially calculated with decorrelated jitter applied.
func ExponentialWithDecorrelatedJitter() func(minDelay, maxDelay time.Duration, attempt int) (delay time.Duration) {
	mutex := &sync.Mutex{}

	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		previous := time.Duration(math.Pow(2, float64(attempt-1)) * float64(minDelay))

		backoff = time.Duration(math.Pow(2, float64(attempt)) * float64(minDelay))

		mutex.Lock()
		jittered := jitter.Decorrelated(minDelay, maxDelay, previous)
		mutex.Unlock()

		backoff += jittered

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}
