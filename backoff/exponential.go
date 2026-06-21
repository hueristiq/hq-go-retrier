package backoff

import (
	"math"
	"time"

	"github.com/hueristiq/hq-lib-retrier-go/jitter"
)

// exponential computes the base exponential backoff for a given attempt.
//
// The base delay grows as minDelay * 2^attempt and is capped at maxDelay. Overflow is guarded by
// checking against math.MaxInt64/2 before doubling and by capping as soon as the next doubling
// would exceed maxDelay. It is the shared core used by all exponential strategies in this package.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 1).
//
// Returns:
//   - base (time.Duration): The exponential base delay, capped at maxDelay. Zero when the inputs
//     are invalid.
//   - ok (bool): False when minDelay or maxDelay is non-positive, or attempt is negative, in which
//     case base is zero and no jitter should be applied.
func exponential(minDelay, maxDelay time.Duration, attempt int) (base time.Duration, ok bool) {
	if minDelay <= 0 || maxDelay <= 0 || attempt < 0 {
		return 0, false
	}

	base = minDelay

	for range attempt {
		if base > math.MaxInt64/2 || base*2 > maxDelay {
			return maxDelay, true
		}

		base *= 2
	}

	if base > maxDelay {
		base = maxDelay
	}

	return base, true
}

// Exponential returns a Backoff function that implements a basic exponential backoff strategy.
//
// This strategy calculates the delay by exponentially increasing the base delay (minDelay) based on
// the retry attempt number, using the formula:
//
//	delay = min(maxDelay, minDelay * 2^attempt)
//
// If minDelay or maxDelay is less than or equal to 0, or if attempt is negative, the function
// returns a zero duration. The delay is capped at maxDelay to ensure reasonable retry intervals.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 1).
//
// Returns:
//   - backoff (Backoff): A function that computes the exponential backoff delay, capped at maxDelay.
func Exponential() Backoff {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		backoff, _ = exponential(minDelay, maxDelay, attempt)

		return
	}
}

// ExponentialWithEqualJitter returns a Backoff function that implements exponential backoff with
// equal jitter to add moderate randomness to retry delays.
//
// The exponential base is computed as min(maxDelay, minDelay * 2^attempt), then equal jitter is
// applied via jitter.Equal, yielding a delay in the range [base/2, base]:
//
//	delay = jitter.Equal(min(maxDelay, minDelay * 2^attempt))
//
// If minDelay or maxDelay is less than or equal to 0, or if attempt is negative, the function
// returns a zero duration.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 1).
//
// Returns:
//   - backoff (Backoff): A function that computes the exponential backoff delay with equal jitter,
//     never exceeding maxDelay.
func ExponentialWithEqualJitter() Backoff {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		base, ok := exponential(minDelay, maxDelay, attempt)
		if !ok {
			return
		}

		backoff = jitter.Equal(base)

		return
	}
}

// ExponentialWithFullJitter returns a Backoff function that implements exponential backoff with
// full jitter to add maximum randomness to retry delays.
//
// The exponential base is computed as min(maxDelay, minDelay * 2^attempt), then full jitter is
// applied via jitter.Full, yielding a delay in the range [0, base):
//
//	delay = jitter.Full(min(maxDelay, minDelay * 2^attempt))
//
// Allowing delays near zero maximally spreads retries, mitigating the "thundering herd" problem.
// If minDelay or maxDelay is less than or equal to 0, or if attempt is negative, the function
// returns a zero duration.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 1).
//
// Returns:
//   - backoff (Backoff): A function that computes the exponential backoff delay with full jitter,
//     never exceeding maxDelay.
func ExponentialWithFullJitter() Backoff {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		base, ok := exponential(minDelay, maxDelay, attempt)
		if !ok {
			return
		}

		backoff = jitter.Full(base)

		return
	}
}

// ExponentialWithDecorrelatedJitter returns a Backoff function that implements exponential backoff
// with decorrelated jitter, reducing correlation between successive retry delays.
//
// The delay is drawn from the range [minDelay, min(maxDelay, previous * 3)], where previous is the
// exponential base of the prior attempt (minDelay * 2^(attempt-1), or minDelay for attempt <= 0):
//
//	delay = jitter.Decorrelated(minDelay, maxDelay, previous)
//
// If minDelay or maxDelay is less than or equal to 0, or if attempt is negative, the function
// returns a zero duration. The delay never exceeds maxDelay.
//
// Parameters:
//   - minDelay (time.Duration): The base (minimum) delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int): The current retry attempt number (typically starting at 1).
//
// Returns:
//   - backoff (Backoff): A function that computes the exponential backoff delay with decorrelated
//     jitter, never exceeding maxDelay.
func ExponentialWithDecorrelatedJitter() Backoff {
	return func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration) {
		_, ok := exponential(minDelay, maxDelay, attempt)
		if !ok {
			return
		}

		previous := minDelay

		if attempt > 0 {
			previous, _ = exponential(minDelay, maxDelay, attempt-1)
		}

		backoff = jitter.Decorrelated(minDelay, maxDelay, previous)

		return
	}
}
