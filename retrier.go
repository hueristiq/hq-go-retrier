package retrier

import (
	"time"

	"go.source.hueristiq.com/retrier/backoff"
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
//   - notifier: A callback function that gets triggered on each retry attempt, providing feedback on errors and backoff duration.
type Configuration struct {
	maxRetries int
	minDelay   time.Duration
	maxDelay   time.Duration
	backoff    backoff.Backoff
	notifier   Notifer
}

// Notifer is a callback function type used to handle notifications during retry attempts.
// This function is invoked on every retry attempt, providing details about the error that
// triggered the retry and the calculated backoff duration before the next attempt.
//
// Parameters:
//   - err: The error encountered in the current retry attempt.
//   - backoff: The duration of backoff calculated before the next retry attempt.
//
// Example:
//
//	func logNotifier(err error, backoff time.Duration) {
//	    fmt.Printf("Retrying after error: %v, backoff: %v\n", err, backoff)
//	}
type Notifer func(err error, backoff time.Duration)

// Option is a function type used to modify the Configuration of the retrier. Options allow
// for the flexible configuration of retry policies by applying user-defined settings.
//
// Parameters:
//   - *Configuration: A pointer to the Configuration struct that allows modification of its fields.
//
// Returns:
//   - Option: A functional option that modifies the Configuration struct, allowing customization of retry behavior.
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

// WithNotifier sets a notifier callback function that gets called on each retry attempt. This function
// allows users to log, monitor, or perform any action upon each retry attempt by providing error details
// and the duration of the backoff period.
//
// Parameters:
//   - notifier: A function of type Notifer that will be called on each retry with the error and backoff duration.
//
// Returns:
//   - Option: A functional option that modifies the Configuration to set the notifier function.
//
// Example:
//
//	retrier.WithNotifier(logNotifier) sets up a notifier that logs each retry attempt.
func WithNotifier(notifier Notifer) Option {
	return func(c *Configuration) {
		c.notifier = notifier
	}
}
