package backoff

import (
	"math"
	"sync"
	"time"

	"github.com/hueristiq/hq-go-retrier/jitter"
)

// Backoff is a function type that calculates the delay between retry attempts.
// It takes a minimum duration, maximum duration, and the current attempt number
// as inputs and returns the calculated delay duration.
type Backoff func(minDelay, maxDelay time.Duration, attempt int) (delay time.Duration)

// Exponential returns a backoff function that implements exponential backoff.
// The delay increases exponentially based on the retry attempt number, with
// the delay capped by the specified maximum duration.
//
// Formula: delay = minDelay * 2^attempt
//
// Parameters:
//   - minDelay: The minimum backoff duration (base duration).
//   - maxDelay: The maximum allowable backoff duration.
//   - attempt: The current retry attempt number.
//
// Returns:
//   - delay: The calculated delay duration, capped at the maximum duration.
func Exponential() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		// Calculate the exponential backoff delay based on the attempt number.
		backoff = time.Duration(math.Pow(2, float64(attempt)) * float64(minDelay))

		// Cap the delay at the maximum value.
		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithEqualJitter returns a backoff function that implements exponential
// backoff with equal jitter. The base delay is calculated using exponential backoff,
// and a random jitter is added, calculated from the midpoint of the delay.
//
// Formula: delay = minDelay * 2^attempt + random(midpoint, delay)
//
// Parameters:
//   - minDelay: The minimum backoff duration (base duration).
//   - maxDelay: The maximum allowable backoff duration.
//   - attempt: The current retry attempt number.
//
// Returns:
//   - delay: The calculated delay with equal jitter applied, capped at the maximum duration.
func ExponentialWithEqualJitter() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	mutex := &sync.Mutex{}

	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		// Calculate the base exponential backoff delay.
		backoff = time.Duration(math.Pow(2, float64(attempt)) * float64(minDelay))

		// Lock the mutex to ensure thread-safe jitter calculation.
		mutex.Lock()
		jittered := jitter.Equal(backoff)
		mutex.Unlock()

		// Add the jitter to the base backoff.
		backoff += jittered

		// Cap the delay at the maximum value.
		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithFullJitter returns a backoff function that implements exponential
// backoff with full jitter. The base delay is calculated using exponential backoff,
// and a random jitter is added to the delay, drawn from a uniform distribution
// in the range [0, delay].
//
// Formula: delay = minDelay * 2^attempt + random(0, delay)
//
// Parameters:
//   - minDelay: The minimum backoff duration (base duration).
//   - maxDelay: The maximum allowable backoff duration.
//   - attempt: The current retry attempt number.
//
// Returns:
//   - delay: The calculated delay with full jitter applied, capped at the maximum duration.
func ExponentialWithFullJitter() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	mutex := &sync.Mutex{}

	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		// Calculate the base exponential backoff delay.
		backoff = time.Duration(math.Pow(2, float64(attempt)) * float64(minDelay))

		// Lock the mutex to ensure thread-safe jitter calculation.
		mutex.Lock()
		jittered := jitter.Full(backoff)
		mutex.Unlock()

		// Add the jitter to the base backoff.
		backoff += jittered

		// Cap the delay at the maximum value.
		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithDecorrelatedJitter returns a backoff function that implements
// exponential backoff with decorrelated jitter. The base delay is calculated using
// exponential backoff, and a random jitter is added that is decorrelated from the
// previous backoff attempt, providing better control over retry intervals.
//
// Formula: delay = minDelay * 2^attempt + random(previous * 3, delay)
//
// Parameters:
//   - minDelay: The minimum backoff duration (base duration).
//   - maxDelay: The maximum allowable backoff duration.
//   - attempt: The current retry attempt number.
//
// Returns:
//   - delay: The calculated delay with decorrelated jitter applied, capped at the maximum duration.
func ExponentialWithDecorrelatedJitter() func(minDelay, maxDelay time.Duration, attempt int) (delay time.Duration) {
	mutex := &sync.Mutex{}

	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		// Calculate the previous backoff delay based on the previous attempt number.
		previous := time.Duration(math.Pow(2, float64(attempt-1)) * float64(minDelay))

		// Calculate the base exponential backoff delay for the current attempt.
		backoff = time.Duration(math.Pow(2, float64(attempt)) * float64(minDelay))

		// Lock the mutex to ensure thread-safe jitter calculation.
		mutex.Lock()
		jittered := jitter.Decorrelated(minDelay, maxDelay, previous)
		mutex.Unlock()

		// Add the decorrelated jitter to the base backoff.
		backoff += jittered

		// Cap the delay at the maximum value.
		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}
