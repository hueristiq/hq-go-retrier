package retrier

import (
	"time"

	"github.com/hueristiq/hq-go-retrier/backoff"
)

// configuration holds the settings for retry operations. These settings determine the behavior
// of the retry mechanism, such as the number of retries, the delay between retries, and the backoff
// strategy to be used.
//
// Fields:
//   - retryMax (int): The maximum number of retry attempts allowed before giving up.
//   - retryWaitMin (time.Duration): The minimum delay between retries, serving as the base delay.
//   - retryWaitMax (time.Duration): The maximum allowable delay between retries.
//   - retryBackoff (backoff.Backoff): A function that calculates the backoff duration based on
//     the current attempt number and the provided delay limits.
//   - notifier (Notifier): A callback function that gets triggered on each retry attempt, allowing
//     for logging or other custom actions based on errors and backoff durations.
type configuration struct {
	retryMax     int
	retryWaitMin time.Duration
	retryWaitMax time.Duration
	retryBackoff backoff.Backoff
	notifier     Notifier
}

// Notifier is a callback function type used to handle notifications during retry attempts.
// This function is invoked on every retry attempt and provides the error that triggered the retry
// along with the calculated backoff duration before the next attempt.
//
// Parameters:
//   - err (error): The error encountered during the current retry attempt.
//   - backoff (time.Duration): The computed delay before the next retry attempt.
type Notifier func(err error, backoff time.Duration)

// Option is a function type used to modify the configuration for the retry mechanism.
// Options allow for flexible and declarative configuration of retry policies by applying
// user-defined settings to a configuration instance.
//
// Parameters:
//   - configuration (*configuration): A pointer to the configuration struct that can be modified..
type Option func(configuration *configuration)

// WithRetryMax returns an Option that sets the maximum number of retry attempts.
//
// When applied, this option limits the number of retries to the specified value. Once the number
// of retry attempts reaches this maximum, the retrier stops further attempts and returns the last error.
//
// Parameters:
//   - retryMax (int): The maximum number of retry attempts.
//
// Returns:
//   - Option: A functional option that updates the retryMax field in the configuration.
func WithRetryMax(retryMax int) (option Option) {
	return func(configuration *configuration) {
		configuration.retryMax = retryMax
	}
}

// WithRetryWaitMin returns an Option that sets the minimum delay between retry attempts.
//
// This option defines the base delay duration from which backoff calculations start, ensuring
// that retries do not occur too rapidly in succession.
//
// Parameters:
//   - retryWaitMin (time.Duration): The minimum delay duration between retries.
//
// Returns:
//   - Option: A functional option that updates the retryWaitMin field in the configuration.
func WithRetryWaitMin(retryWaitMin time.Duration) (option Option) {
	return func(configuration *configuration) {
		configuration.retryWaitMin = retryWaitMin
	}
}

// WithRetryWaitMax returns an Option that sets the maximum delay between retry attempts.
//
// This option imposes an upper bound on the delay between retries, ensuring that the backoff
// duration does not grow unbounded even with strategies like exponential backoff.
//
// Parameters:
//   - retryWaitMax (time.Duration): The maximum allowable delay duration between retries.
//
// Returns:
//   - Option: A functional option that updates the retryWaitMax field in the configuration.
func WithRetryWaitMax(retryWaitMax time.Duration) (option Option) {
	return func(configuration *configuration) {
		configuration.retryWaitMax = retryWaitMax
	}
}

// WithRetryBackoff returns an Option that sets the backoff strategy used to compute the delay
// between retry attempts. The backoff strategy is a function (of type backoff.Backoff) that calculates
// the delay based on the current retry attempt number and the defined minimum and maximum delays.
// This allows for various backoff algorithms (such as exponential backoff with jitter) to be applied.
//
// Parameters:
//   - retryBackoff (backoff.Backoff): A function that defines the backoff strategy.
//
// Returns:
//   - option (Option): A functional option that updates the retryBackoff field in the configuration.
func WithRetryBackoff(retryBackoff backoff.Backoff) (option Option) {
	return func(configuration *configuration) {
		configuration.retryBackoff = retryBackoff
	}
}

// WithNotifier returns an Option that sets a notifier callback function for retry attempts.
// The notifier function is invoked on each retry attempt, and it receives the error that caused
// the retry along with the computed backoff duration. This is useful for logging, monitoring,
// or triggering other side effects on retries.
//
// Parameters:
//   - notifier (Notifier): A callback function of type Notifier to be called on each retry.
//
// Returns:
//   - option (Option): A functional option that updates the notifier field in the configuration.
func WithNotifier(notifier Notifier) (option Option) {
	return func(configuration *configuration) {
		configuration.notifier = notifier
	}
}
