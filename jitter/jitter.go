package jitter

import (
	"crypto/rand"
	"math/big"
	"time"
)

// Equal applies an equal jitter strategy to the provided backoff duration.
// In this strategy, the jittered duration is calculated as the midpoint of the
// backoff duration plus a random value uniformly selected from the range [0, midpoint).
// This reduces randomness compared to full jitter but still ensures some variation.
//
// Parameters:
//   - backoff: The original backoff duration before applying jitter.
//
// Returns:
//   - jitter: A duration that is the midpoint of the backoff plus a random value between 0 and the midpoint.
func Equal(backoff time.Duration) (jitter time.Duration) {
	// Calculate the midpoint of the backoff duration.
	midpoint := backoff / 2

	// Add a random duration between 0 and the midpoint to the midpoint.
	jitter = midpoint + getRandomDuration(midpoint)

	return
}

// Full applies a full jitter strategy to the provided backoff duration.
// In this strategy, the jittered duration is a random value uniformly selected
// from the range [0, backoff). This method introduces maximum randomness.
//
// Parameters:
//   - backoff: The original backoff duration before applying jitter.
//
// Returns:
//   - jitter: A random duration between 0 and the provided backoff duration.
func Full(backoff time.Duration) (jitter time.Duration) {
	// Generate a random duration between 0 and the backoff duration.
	jitter = getRandomDuration(backoff)

	return
}

// Decorrelated applies a decorrelated jitter strategy, which uses the previous
// backoff duration to calculate the new jittered duration. This strategy
// prevents exponential backoff from growing too large, by bounding it within
// the provided minimum and maximum values.
//
// Parameters:
//   - min: The minimum backoff duration.
//   - max: The maximum backoff duration.
//   - previous: The previous backoff duration, which influences the next jittered value.
//
// Returns:
//   - jitter: A jittered duration that is random but decorrelated from the previous one.
//     The returned duration is within the range [minDelay, maxDelay].
func Decorrelated(minDelay, maxDelay, previous time.Duration) (jitter time.Duration) {
	// If this is the first call, use the minimum duration as the previous value.
	if previous == 0 {
		previous = minDelay
	}

	// Generate a random duration within the range [minDelay, previous*3].
	jitter = getRandomDuration(previous * 3)
	jitter += minDelay

	// Ensure that the jitter does not exceed the maximum duration.
	if jitter > maxDelay {
		jitter = maxDelay
	}

	return
}

// getRandomDuration returns a random time.Duration value between 0 and maxDuration.
// It uses a cryptographically secure random number generator (CSPRNG) to produce
// the random value, ensuring a high degree of unpredictability.
//
// Parameters:
//   - maxDuration: The maxDurationimum possible duration for the random value. Must be greater than 0.
//
// Returns:
//   - A random duration in the range [0, maxDuration]. If maxDuration is less than or equal to 0, returns 0.
func getRandomDuration(maxDuration time.Duration) (duration time.Duration) {
	// Return 0 if the maxDurationimum value is invalid or non-positive.
	if maxDuration <= 0 {
		return 0
	}

	// Generate a cryptographically secure random integer between 0 and maxDuration.
	n, err := rand.Int(rand.Reader, big.NewInt(int64(maxDuration)))
	if err != nil {
		// If an error occurs during random number generation, return the maxDurationimum value.
		return maxDuration
	}

	// Convert the random integer to a time.Duration and return it.
	duration = time.Duration(n.Int64())

	return
}
