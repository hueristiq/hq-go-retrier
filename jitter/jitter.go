package jitter

import (
	"crypto/rand"
	"math/big"
	"time"
)

// Equal applies an equal jitter strategy to the provided backoff duration.
// This strategy introduces moderate randomness by taking the midpoint of the original
// backoff duration and adding a random offset between 0 and that midpoint.
//
// The calculation is as follows:
//
//	jitter = (backoff / 2) + random(0, backoff / 2)
//
// This ensures that the resulting delay is between backoff/2 and backoff.
// Equal jitter is useful when you want a relatively stable retry interval with some variation,
// but without the extreme fluctuations that full jitter can introduce.
//
// Parameters:
//   - backoff (time.Duration): The original backoff duration (base delay) to which jitter is applied.
//
// Returns:
//   - jitter (time.Duration): The resulting delay after applying equal jitter, which will be a value
//     between backoff/2 and backoff.
//
// Example:
//
//	baseDelay := 10 * time.Second
//	jitteredDelay := Equal(baseDelay)
//	// jitteredDelay will be between 5 seconds and 10 seconds.
func Equal(backoff time.Duration) (jitter time.Duration) {
	midpoint := backoff / 2

	jitter = midpoint + getRandomDuration(midpoint)

	return
}

// Full applies a full jitter strategy to the provided backoff duration.
// In this strategy, the delay is completely randomized by selecting a random value
// uniformly between 0 and the original backoff duration.
//
// This approach introduces maximum randomness, which can be effective in spreading out
// retry attempts widely when multiple clients or operations are retrying simultaneously.
//
// Parameters:
//   - backoff (time.Duration): The base backoff duration to be randomized.
//
// Returns:
//   - jitter (time.Duration): A random delay between 0 and backoff.
//
// Example:
//
//	baseDelay := 10 * time.Second
//	jitteredDelay := Full(baseDelay)
//	// jitteredDelay will be a random duration between 0 and 10 seconds.
func Full(backoff time.Duration) (jitter time.Duration) {
	jitter = getRandomDuration(backoff)

	return
}

// Decorrelated applies a decorrelated jitter strategy to the backoff duration.
// This strategy not only increases the delay exponentially but also introduces randomness
// in a way that is influenced by the previous backoff duration. This helps in avoiding
// exponential growth that could otherwise lead to excessively long delays.
//
// The algorithm works as follows:
//   - If previous is zero, it defaults to minDelay.
//   - A random duration is generated in the range [0, previous * 3).
//   - This random duration is added to minDelay.
//   - The resulting delay is capped by maxDelay to ensure it does not exceed the allowed maximum.
//
// This method allows each backoff duration to be decorrelated from its predecessor while
// ensuring that the delay remains bounded between minDelay and maxDelay.
//
// Parameters:
//   - minDelay (time.Duration): The minimum backoff duration (base delay).
//   - maxDelay (time.Duration): The maximum allowable backoff duration.
//   - previous (time.Duration): The previous backoff duration, which influences the amount of jitter applied.
//
// Returns:
//   - jitter (time.Duration): The new backoff duration after applying decorrelated jitter, guaranteed
//     to be between minDelay and maxDelay.
//
// Example:
//
//	minDelay := 1 * time.Second
//	maxDelay := 30 * time.Second
//	previousDelay := 5 * time.Second
//	jitteredDelay := Decorrelated(minDelay, maxDelay, previousDelay)
//	// jitteredDelay will be a decorrelated value influenced by previousDelay and capped at maxDelay.
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

// getRandomDuration returns a random time.Duration value between 0 and the provided maximum duration.
// This function uses a cryptographically secure random number generator (CSPRNG) via crypto/rand to ensure
// that the random values are highly unpredictable.
//
// Parameters:
//   - maxDuration (time.Duration): The upper bound for the random duration. Must be a positive value.
//
// Returns:
//   - duration (time.Duration): A random duration in the range [0, maxDuration). If maxDuration is less than
//     or equal to zero, the function returns 0.
//
// Example:
//
//	randomDelay := getRandomDuration(10 * time.Second)
//	// randomDelay will be a random duration between 0 and 10 seconds.
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
