package backoff

import (
	"time"
)

// Backoff is a function type that calculates the delay duration between successive retry attempts.
//
// It is designed to be used as an abstraction for various retry strategies (e.g., constant, linear,
// exponential, or jittered backoff) by defining a common function signature. Implementations of this
// function type determine how long a caller should wait before attempting a retry after a failure.
//
// Parameters:
//   - minDelay (time.Duration): The minimum allowable delay duration. This value represents the
//     smallest amount of time that should be waited between retries, regardless of the retry attempt.
//   - maxDelay (time.Duration): The maximum allowable delay duration. This value acts as an upper
//     bound on the delay, ensuring that the wait time does not exceed a specified limit even as the
//     number of attempts increases.
//   - attempt (int): The current retry attempt number. This value is typically a positive integer
//     (starting at 1) that indicates how many times a retry has been attempted. Implementations can use
//     this parameter to increase the delay in subsequent attempts (e.g., doubling the delay on each retry).
//
// Returns:
//   - delay (time.Duration): The calculated delay duration to wait before the next retry attempt.
//     This delay is typically derived by applying a strategy (such as exponential growth) to the input
//     parameters, while ensuring that the returned value is not less than minDelay or greater than maxDelay.
type Backoff func(minDelay, maxDelay time.Duration, attempt int) (delay time.Duration)
