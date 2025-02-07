package jitter

import (
	"crypto/rand"
	"math/big"
	"time"
)

// Equal applies an equal jitter strategy to the provided backoff duration.
// This method ensures moderate randomness by adding a jitter value that is
// calculated as a random number within half of the original backoff time.
//
// Equal jitter provides a balance between fixed backoff and full jitter by
// ensuring that the backoff duration does not deviate too much, reducing
// the likelihood of excessively short or long delays. This approach is
// useful for scenarios where consistent retry intervals with a bit of
// variation are preferred.
//
// Arguments:
//   - backoff: The original backoff duration to which jitter will be applied.
//     This represents the base amount of time to wait before retrying
//     an operation.
//
// Returns:
//   - jitter: The resulting backoff duration after applying equal jitter.
//     It will be the midpoint of the original backoff plus a random
//     value between 0 and the midpoint.
//
// Example:
//
//	backoff := 10 * time.Second
//	jitteredBackoff := jitter.Equal(backoff)
//	// jitteredBackoff will be somewhere between 5 seconds and 10 seconds.
func Equal(backoff time.Duration) (jitter time.Duration) {
	midpoint := backoff / 2

	jitter = midpoint + getRandomDuration(midpoint)

	return
}

// Full applies a full jitter strategy to the provided backoff duration.
// In this strategy, a random value is selected from the entire range between
// 0 and the original backoff, providing maximum randomness to the retry delay.
//
// Full jitter is useful when you want to distribute the retry attempts more
// uniformly across a wide range of possible durations, avoiding scenarios
// where retries happen at consistent intervals, which might overload a
// system under high contention.
//
// Arguments:
//   - backoff: The base backoff duration to be randomized.
//
// Returns:
//   - jitter: A completely random backoff duration between 0 and the original
//     backoff value.
//
// Example:
//
//	backoff := 10 * time.Second
//	jitteredBackoff := jitter.Full(backoff)
//	// jitteredBackoff will be somewhere between 0 and 10 seconds.
func Full(backoff time.Duration) (jitter time.Duration) {
	jitter = getRandomDuration(backoff)

	return
}

// Decorrelated applies a decorrelated jitter strategy to the backoff duration.
// This method calculates the jittered duration using a random value from a
// range that is influenced by the previous backoff value. This approach
// prevents exponential growth of backoff intervals, which could otherwise
// increase uncontrollably in some retry scenarios.
//
// The decorrelated jitter strategy ensures that the backoff is randomized
// but bounded, making it effective for systems where exponential backoff
// needs to be capped to avoid overly long retry intervals.
//
// Arguments:
//   - minDelay: The minimum delay duration for the backoff.
//   - maxDelay: The maximum allowable delay duration for the backoff.
//   - previous: The previous backoff duration, used to calculate the new
//     jittered duration.
//
// Returns:
//   - jitter: A decorrelated jittered duration that is within the range of
//     [minDelay, maxDelay]. The next backoff will be influenced by
//     the previous one but bounded to avoid excessive delays.
//
// Example:
//
//	minDelay := 1 * time.Second
//	maxDelay := 30 * time.Second
//	previous := 5 * time.Second
//	jitteredBackoff := jitter.Decorrelated(minDelay, maxDelay, previous)
//	// jitteredBackoff will be somewhere between minDelay and maxDelay,
//	// bounded by the previous backoff value.
func Decorrelated(minDelay, maxDelay, previous time.Duration) (jitter time.Duration) {
	if previous == 0 {
		previous = minDelay
	}

	jitter = getRandomDuration(previous * 3)

	jitter += minDelay

	if jitter > maxDelay {
		jitter = maxDelay
	}

	return
}

// getRandomDuration returns a random time.Duration value between 0 and the
// provided maximum duration. This function uses a cryptographically secure
// random number generator (CSPRNG) to ensure that the random values are
// highly unpredictable.
//
// By using a CSPRNG, this method guarantees stronger randomness compared
// to traditional pseudo-random number generators, which is particularly
// useful in security-sensitive applications.
//
// Arguments:
//   - maxDuration: The maximum duration from which to select a random value.
//     This must be a positive value greater than zero.
//
// Returns:
//   - duration: A random time.Duration value between 0 and maxDuration. If
//     maxDuration is less than or equal to 0, the function returns
//     a duration of 0.
//
// Example:
//
//	randomDuration := getRandomDuration(10 * time.Second)
//	// randomDuration will be a random time.Duration between 0 and 10 seconds.
func getRandomDuration(maxDuration time.Duration) (duration time.Duration) {
	if maxDuration <= 0 {
		return 0
	}

	n, err := rand.Int(rand.Reader, big.NewInt(int64(maxDuration)))
	if err != nil {
		return maxDuration
	}

	duration = time.Duration(n.Int64())

	return
}
