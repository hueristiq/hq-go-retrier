package backoff

import (
	"time"
)

// Backoff is a function type that calculates the delay between retry attempts.
// It takes a minimum duration, maximum duration, and the current retry attempt number as inputs,
// and returns the calculated delay duration. Each retry strategy returns a function of this type.
//
// Arguments:
//   - minDelay (time.Duration): The minimum allowable delay duration.
//   - maxDelay (time.Duration): The maximum allowable delay duration.
//   - attempt (int):  The current retry attempt number.
//
// Returns:
//   - delay (time.Duration): The delay duration to wait before the next retry.
type Backoff func(minDelay, maxDelay time.Duration, attempt int) (delay time.Duration)
