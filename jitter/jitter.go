package jitter

import (
	"crypto/rand"
	"math/big"
	"time"
)

// Equal calculates a jitter duration using the equal jitter strategy based on the provided backoff.
//
// This strategy introduces moderate randomness by calculating a midpoint of the backoff duration
// and adding a random offset between 0 and that midpoint. The resulting jitter duration is guaranteed
// to be in the range [backoff/2, backoff], balancing predictability and randomization to avoid
// synchronized retry spikes in distributed systems.
//
//	jitter = (backoff / 2) + random(0, backoff / 2)
//
// Parameters:
//   - backoff (time.Duration): The base backoff duration to which jitter is applied.
//
// Returns:
//   - jitter (time.Duration): The calculated jitter duration, always in [backoff/2, backoff] for positive backoff.
//     Returns 0 if backoff is 0 or negative.
func Equal(backoff time.Duration) (jitter time.Duration) {
	jitter = 0

	if backoff <= 0 {
		return
	}

	midpoint := backoff / 2

	jitter = midpoint + getRandomDuration(midpoint)

	return
}

// Full calculates a jitter duration using the full jitter strategy based on the provided backoff.
//
// This strategy generates a fully randomized delay between 0 and the backoff duration, providing
// maximum randomness to spread out retry attempts and mitigate the "thundering herd" problem.
//
//	jitter = random(0, backoff)
//
// Parameters:
//   - backoff (time.Duration): The base backoff duration to which jitter is applied.
//
// Returns:
//   - jitter (time.Duration): The calculated jitter duration, in [0, backoff] for positive backoff.
//     Returns 0 if backoff is 0 or negative.
func Full(backoff time.Duration) (jitter time.Duration) {
	jitter = 0

	if backoff <= 0 {
		return
	}

	jitter = getRandomDuration(backoff)

	return
}

// Decorrelated calculates a jitter duration using the decorrelated jitter strategy,
// incorporating the previous backoff to reduce correlation between successive retries.
//
// This strategy introduces randomness influenced by the previous backoff duration, while enforcing
// minimum and maximum delay bounds. It is designed to prevent exponential backoff growth from
// becoming excessive, while still providing randomness to avoid synchronized retries.
// The jitter is calculated as follows:
//   - If previous is 0, it is set to minDelay.
//   - A random duration is generated in the range [0, previous * 3).
//   - The random duration is added to minDelay.
//   - The result is capped at maxDelay to ensure it does not exceed the maximum allowed delay.
//
// Parameters:
//   - minDelay (time.Duration): The minimum allowable jitter duration.
//   - maxDelay (time.Duration): The maximum allowable jitter duration.
//   - previous (time.Duration): The previous backoff duration, influencing the jitter range.
//     If `0`, defaults to `minDelay`.
//
// Returns:
//   - jitter (time.Duration): The calculated jitter duration, in [minDelay, maxDelay].
//     Returns minDelay if previous * 3 is less than minDelay, or maxDelay if the calculated jitter exceeds it.
func Decorrelated(minDelay, maxDelay, previous time.Duration) (jitter time.Duration) {
	jitter = 0

	if minDelay < 0 || maxDelay < 0 || minDelay > maxDelay {
		return
	}

	if previous <= 0 {
		previous = minDelay
	}

	jitter = getRandomDuration(previous * 3)

	jitter += minDelay

	if jitter > maxDelay {
		jitter = maxDelay
	}

	return
}

// getRandomDuration generates a cryptographically secure random duration between 0 and maxDuration.
//
// It uses the crypto/rand package to ensure high-quality randomness, suitable for jitter calculations
// in retry strategies where unpredictability is critical.
//
// Parameters:
//   - maxDuration (time.Duration): The upper bound for the random duration. Must be positive.
//     If 0 or negative, returns 0, as no meaningful random duration can be generated.
//
// Returns:
//   - duration (time.Duration): A random duration in the range [`0`, `maxDuration`).
//     Returns `0` if `maxDuration <= 0`.
//     Returns `maxDuration` if the random number generation fails (fallback for robustness).
func getRandomDuration(maxDuration time.Duration) (duration time.Duration) {
	duration = 0

	if maxDuration <= 0 {
		return
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(maxDuration)))
	if err != nil {
		duration = maxDuration

		return
	}

	duration = time.Duration(n.Int64())

	return
}
