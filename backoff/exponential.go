package backoff

import (
	"time"

	"go.source.hueristiq.com/retrier/jitter"
)

// Exponential returns a backoff function that implements basic exponential backoff.
// In this strategy, the delay increases exponentially with each retry attempt.
// The delay is calculated using the formula:
//
//	delay = minDelay * 2^(attempt)
//
// The calculated delay is then capped by maxDelay, ensuring that the delay never exceeds
// a specified maximum.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 0 or 1).
//
// Returns:
//   - backoff (time.Duration): The computed delay duration to wait before the next retry.
//
// Example:
//
//	backoffFunc := Exponential()
//	delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	// For attempt 3, the base delay is calculated as 1s * 2^3 = 8 seconds,
//	// but if the calculated delay exceeds maxDelay, it is capped at maxDelay.
func Exponential() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = minDelay << attempt

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithEqualJitter returns a backoff function that implements exponential backoff with equal jitter.
// In this strategy, the delay is first calculated exponentially and then a random jitter value is added.
// The jitter is computed using an "equal" jitter function, which typically adds a random value around the midpoint
// of the calculated delay. The formula can be described as:
//
//	delay = (minDelay * 2^(attempt)) + jitter.Equal(delay)
//
// As with basic exponential backoff, the final delay is capped at maxDelay.
//
// Parameters:
//   - minDelay (time.Duration): The base delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number.
//
// Returns:
//   - backoff (time.Duration): The computed delay with equal jitter applied.
//
// Example:
//
//	backoffFunc := ExponentialWithEqualJitter()
//	delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	// The base delay is 8 seconds (1s * 2^3); jitter.Equal adds a moderate random delay,
//	// and the total is capped at 30 seconds if necessary.
func ExponentialWithEqualJitter() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = minDelay << attempt

		backoff += jitter.Equal(backoff)

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithFullJitter returns a backoff function that implements exponential backoff with full jitter.
// In this strategy, after calculating the exponential delay, a random jitter is applied such that the additional
// delay is uniformly distributed between 0 and the computed delay. The formula is:
//
//	delay = (minDelay * 2^(attempt)) + jitter.Full(delay)
//
// This full jitter approach often results in a more randomized delay to help avoid the "thundering herd" effect,
// and the final delay is capped at maxDelay.
//
// Parameters:
//   - minDelay (time.Duration): The base delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number.
//
// Returns:
//   - backoff (time.Duration): The computed delay with full jitter applied.
//
// Example:
//
//	backoffFunc := ExponentialWithFullJitter()
//	delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	// The base delay is calculated as 8 seconds (1s * 2^3), then a random value between 0 and 8 seconds is added,
//	// and the result is capped at 30 seconds.
func ExponentialWithFullJitter() func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = minDelay << attempt

		backoff += jitter.Full(backoff)

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithDecorrelatedJitter returns a backoff function that implements exponential backoff with decorrelated jitter.
// This strategy not only increases the delay exponentially, but also introduces jitter in a way that reduces correlation
// between successive retries. It takes into account the previous backoff duration when computing the jitter.
// The general idea is:
//
//	delay = (minDelay * 2^(attempt)) + jitter.Decorrelated(minDelay, maxDelay, previous)
//
// where 'previous' is the backoff value from the previous attempt (or minDelay if no previous attempt exists).
// If the provided attempt is less than 0, the function returns minDelay.
//
// The final computed delay is capped by maxDelay.
//
// Parameters:
//   - minDelay (time.Duration): The base delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number.
//   - If attempt < 0, the function returns minDelay immediately.
//   - For attempt > 0, the previous delay is computed as minDelay * 2^(attempt - 1).
//
// Returns:
//   - delay (time.Duration): The computed delay with decorrelated jitter applied.
//
// Example:
//
//	backoffFunc := ExponentialWithDecorrelatedJitter()
//	delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	// The delay is calculated with exponential growth, plus decorrelated jitter based on the previous delay,
//	// and is capped at 30 seconds.
func ExponentialWithDecorrelatedJitter() func(minDelay, maxDelay time.Duration, attempt int) (delay time.Duration) {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		if attempt < 0 {
			return minDelay
		}

		previous := minDelay

		if attempt > 0 {
			previous = minDelay << (attempt - 1)
		}

		backoff = minDelay << attempt

		backoff += jitter.Decorrelated(minDelay, maxDelay, previous)

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}
