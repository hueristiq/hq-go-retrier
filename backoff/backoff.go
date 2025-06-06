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
// Parameters:
//   - minDelay (time.Duration): The minimum allowable delay duration. Ensures that the returned
//     backoff duration is at least this value, preventing excessively short delays that could
//     lead to rapid retry attempts.
//   - maxDelay (time.Duration): The maximum allowable delay duration. Caps the returned backoff
//     duration to prevent excessively long delays, ensuring retries occur within a reasonable
//     timeframe.
//   - attempt (int): Typically, maxDelay should be greater than or equal to minDelay.
//   - attempt (int): The current retry attempt number, typically starting at 1 for the first retry.
//     Implementations use this value to adjust the delay (e.g., increasing it for subsequent retries
//     in exponential backoff).
//
// Returns:
//   - backoff (time.Duration): The calculated delay duration to wait before the next retry attempt.
//     The duration is computed based on the implemented strategy (e.g., constant, linear, or
//     exponential) and is guaranteed to be in the range [minDelay, maxDelay].
type Backoff func(minDelay, maxDelay time.Duration, attempt int) (backoff time.Duration)
