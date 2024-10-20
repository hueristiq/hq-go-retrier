package retrier

import (
	"time"

	"github.com/hueristiq/hq-go-retrier/backoff"
)

// Configuration holds the settings for retry operations. These settings determine the behavior
// of the retry mechanism, such as the number of retries, delay between retries, and the backoff
// strategy to be used.
//
// Fields:
//   - maxRetries: The maximum number of retry attempts allowed before giving up.
//   - minDelay: The minimum delay between retries.
//   - maxDelay: The maximum allowable delay between retries.
//   - backoff: A function that calculates the backoff duration based on retry attempt number and delay limits.
type Configuration struct {
	maxRetries int             // Maximum number of retry attempts.
	minDelay   time.Duration   // Minimum delay between retry attempts.
	maxDelay   time.Duration   // Maximum delay between retry attempts.
	backoff    backoff.Backoff // Backoff strategy used to calculate delay between attempts.
}

// Option is a function type used to modify the Configuration of the retrier. Options allow
// for the flexible configuration of retry policies by applying user-defined settings.
type Option func(*Configuration)

// WithMaxRetries sets the maximum number of retries for the retry mechanism. When the specified
// number of retries is reached, the operation will stop, and the last error will be returned.
//
// Parameters:
//   - retries: The maximum number of retry attempts.
//
// Returns:
//   - Option: A functional option that modifies the Configuration to set the maxRetries field.
//
// Example:
//
//	retrier.WithMaxRetries(5) sets the retry policy to attempt a maximum of 5 retries.
func WithMaxRetries(retries int) Option {
	return func(c *Configuration) {
		c.maxRetries = retries
	}
}

// WithMaxDelay sets the maximum allowable delay between retry attempts. This option ensures that
// the delay between retries never exceeds the specified maximum, even with exponential backoff.
//
// Parameters:
//   - delay: The maximum delay duration between retries.
//
// Returns:
//   - Option: A functional option that modifies the Configuration to set the maxDelay field.
//
// Example:
//
//	retrier.WithMaxDelay(2 * time.Second) ensures that delays between retries do not exceed 2 seconds.
func WithMaxDelay(delay time.Duration) Option {
	return func(c *Configuration) {
		c.maxDelay = delay
	}
}

// WithMinDelay sets the minimum delay between retry attempts. This is the base duration from which
// the delay calculations start, and it ensures that retries do not occur too quickly in rapid succession.
//
// Parameters:
//   - delay: The minimum delay duration between retries.
//
// Returns:
//   - Option: A functional option that modifies the Configuration to set the minDelay field.
//
// Example:
//
//	retrier.WithMinDelay(100 * time.Millisecond) ensures that retries wait at least 100ms before retrying.
func WithMinDelay(delay time.Duration) Option {
	return func(c *Configuration) {
		c.minDelay = delay
	}
}

// WithBackoff sets the backoff strategy used to calculate the delay between retry attempts. The backoff
// strategy determines how the delay grows between retries, and can be customized to use strategies such as
// exponential backoff with jitter.
//
// Parameters:
//   - strategy: A backoff function that defines the backoff strategy.
//
// Returns:
//   - Option: A functional option that modifies the Configuration to set the backoff strategy.
//
// Example:
//
//	retrier.WithBackoff(backoff.ExponentialWithFullJitter()) configures the retrier to use exponential backoff with full jitter.
func WithBackoff(strategy backoff.Backoff) Option {
	return func(c *Configuration) {
		c.backoff = strategy
	}
}
