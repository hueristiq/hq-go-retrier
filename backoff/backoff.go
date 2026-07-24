package backoff

import (
	"time"
)

// Backoff is a function type that calculates the delay duration between successive retry attempts.
//
// It serves as an abstraction for implementing various retry strategies (e.g., constant, linear,
// exponential, or jittered backoff) by defining a standardized function signature. Implementations
// of this type determine the wait time before a retry attempt, allowing for flexible and customizable
// retry logic in fault-tolerant systems.
//
// Implementations must be safe for concurrent use by multiple goroutines, since a single Backoff
// value may be shared across independent retry loops.
//
// Parameters:
//   - minDelay (time.Duration): The minimum allowable delay duration. Strategies use it as the
//     base (or floor) of the delay, preventing excessively short delays that could lead to
//     rapid retry attempts. Jittered strategies may dip below it; see the Returns section.
//   - maxDelay (time.Duration): The maximum allowable delay duration. Caps the returned backoff
//     duration to prevent excessively long delays, ensuring retries occur within a reasonable
//     timeframe. Typically, maxDelay should be greater than or equal to minDelay.
//   - attempt (int): The current retry attempt number, typically starting at 1 for the first retry.
//     Implementations use this value to adjust the delay (e.g., increasing it for subsequent retries
//     in exponential backoff).
//
// Returns:
//   - backoff (time.Duration): The calculated delay duration to wait before the next retry attempt,
//     computed based on the implemented strategy (e.g., constant, linear, or exponential). For
//     valid input the delay lies in [0, maxDelay]; jittered strategies may legitimately return
//     less than minDelay. Implementations return a zero duration for invalid input — non-positive
//     bounds, minDelay greater than maxDelay, or a negative attempt.
type Backoff func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration)
