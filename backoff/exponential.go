package backoff

import (
	"math"
	"time"

	"github.com/hueristiq/hq-go-retrier/jitter"
)

// Exponential returns a Backoff function that implements a basic exponential backoff strategy.
//
// This strategy calculates the delay by exponentially increasing the base delay (minDelay) based on
// the retry attempt number, using the formula:
//
//	delay = minDelay * 2^attempt
//
// If minDelay or maxDelay is less than or equal to 0, or if attempt is negative, the function
// returns a zero duration. For attempt < 1, it returns minDelay (no exponential increase).
// The delay is capped at maxDelay to ensure reasonable retry intervals.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 0 or 1).
//
// Returns:
//   - backoff (Backoff): A function that computes the exponential backoff delay, capped at maxDelay.
func Exponential() Backoff {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = 0

		if minDelay <= 0 || maxDelay <= 0 || attempt < 0 {
			return
		}

		backoff = minDelay

		if maxDelay > minDelay && attempt < 1 {
			return
		}

		for range attempt {
			if backoff > math.MaxInt64/2 {
				backoff = maxDelay

				return
			}

			backoff *= 2
		}

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithEqualJitter returns a Backoff function that implements exponential backoff with
// equal jitter to add moderate randomness to retry delays.
//
// The delay is calculated as minDelay * 2^attempt, then augmented with equal jitter from the jitter.Equal
// function, which adds a random duration in the range [delay/2, delay]. The formula is:
//
//	delay = (minDelay * 2^attempt) + jitter.Equal(delay)
//
// The final delay is capped at maxDelay. If minDelay or maxDelay is less than or equal to 0, or if
// attempt is negative, the function returns a zero duration. For attempt < 1, it returns minDelay
// plus equal jitter.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 0 or 1).
//
// Returns:
//   - backoff (Backoff): A function that computes the exponential backoff delay with equal jitter,
//     capped at maxDelay.
func ExponentialWithEqualJitter() Backoff {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = 0

		if minDelay <= 0 || maxDelay <= 0 || attempt < 0 {
			return
		}

		backoff = minDelay

		if maxDelay > minDelay && attempt < 1 {
			return
		}

		for range attempt {
			if backoff > math.MaxInt64/2 {
				backoff = maxDelay

				return
			}

			backoff *= 2
		}

		backoff += jitter.Equal(backoff)

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithFullJitter returns a Backoff function that implements exponential backoff with
// full jitter to add maximum randomness to retry delays.
//
// The delay is calculated as minDelay * 2^attempt, then augmented with full jitter from the jitter.Full
// function, which adds a random duration in the range [0, delay]. The formula is:
//
//	delay = (minDelay * 2^attempt) + jitter.Full(delay)
//
// The final delay is capped at maxDelay. If minDelay or maxDelay is less than or equal to 0, or if
// attempt is negative, the function returns a zero duration. For attempt < 1, it returns minDelay
// plus full jitter.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 0 or 1).
//
// Returns:
//   - backoff (Backoff): A function that computes the exponential backoff delay with full jitter,
//     capped at maxDelay.
func ExponentialWithFullJitter() Backoff {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = 0

		if minDelay <= 0 || maxDelay <= 0 || attempt < 0 {
			return
		}

		backoff = minDelay

		if maxDelay > minDelay && attempt < 1 {
			return
		}

		for range attempt {
			if backoff > math.MaxInt64/2 {
				backoff = maxDelay

				return
			}

			backoff *= 2
		}

		backoff += jitter.Full(backoff)

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}

// ExponentialWithDecorrelatedJitter returns a Backoff function that implements exponential backoff
// with decorrelated jitter, reducing correlation between successive retry delays.
//
// The delay is calculated as minDelay * 2^attempt, then augmented with decorrelated jitter from the
// jitter.Decorrelated function, which computes a random duration in the range [minDelay, previous * 3],
// where previous is the delay of the previous attempt (minDelay * 2^(attempt-1)) or minDelay for
// attempt = 0. The formula is:
//
//	delay = (minDelay * 2^attempt) + jitter.Decorrelated(minDelay, maxDelay, previous)
//
// The final delay is capped at maxDelay. If minDelay or maxDelay is less than or equal to 0, or if
// attempt is negative, the function returns a zero duration. For attempt < 1, it returns minDelay
// plus decorrelated jitter.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 0 or 1).
//
// Returns:
//   - backoff (Backoff): A function that computes the exponential backoff delay with decorrelated
//     jitter, capped at maxDelay.
func ExponentialWithDecorrelatedJitter() Backoff {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff = 0

		if minDelay <= 0 || maxDelay <= 0 || attempt < 0 {
			return
		}

		backoff = minDelay

		if maxDelay > minDelay && attempt < 1 {
			return
		}

		for range attempt {
			if backoff > math.MaxInt64/2 {
				backoff = maxDelay

				return
			}

			backoff *= 2
		}

		previous := minDelay

		if attempt > 0 {
			for range attempt - 1 {
				if previous > math.MaxInt64/2 {
					previous = maxDelay

					break
				}

				previous *= 2
			}
		}

		backoff += jitter.Decorrelated(minDelay, maxDelay, previous)

		if backoff > maxDelay {
			backoff = maxDelay
		}

		return
	}
}
