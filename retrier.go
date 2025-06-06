package retrier

import (
	"context"
	"time"

	"github.com/hueristiq/hq-go-retrier/backoff"
)

// configuration holds the settings for retry operations, defining the behavior of the retry
// mechanism.
//
// Fields:
//   - retryMax (int): The maximum number of retry attempts allowed before giving up.
//   - retryWaitMin (time.Duration): The minimum allowable delay between retry attempts, serving as the
//     base delay for backoff calculations.
//   - retryWaitMax (time.Duration): The maximum allowable delay between retry attempts, capping
//     the backoff duration.
//   - retryBackoff (backoff.Backoff): A function that calculates the backoff duration based on
//     the current attempt number, retryWaitMin, and retryWaitMax.
//   - notifier (Notifier): A callback function invoked on each retry attempt, receiving the error
//     that triggered the retry and the computed backoff duration.
type configuration struct {
	retryMax     int
	retryWaitMin time.Duration
	retryWaitMax time.Duration
	retryBackoff backoff.Backoff
	notifier     Notifier
}

// Notifier is a callback function type used to handle notifications during retry attempts.
//
// It is invoked after each failed retry attempt, providing the error that caused the retry and
// the computed backoff duration before the next attempt. This allows for custom logging,
// monitoring, or other side effects during retries.
//
// Parameters:
//   - err (error): The error encountered during the current retry attempt. Will not be nil.
//   - backoff (time.Duration): The computed delay duration before the next retry attempt.
type Notifier func(err error, backoff time.Duration)

// Option is a function type used to modify the retry configuration in a declarative manner.
//
// It allows users to customize retry behavior by setting fields in the configuration struct,
// such as the maximum number of retries, delay bounds, backoff strategy, or notifier callback.
// Multiple options can be combined to create a tailored retry policy.
//
// Parameters:
//   - configuration (*configuration): A pointer to the configuration struct to be modified.
type Option func(configuration *configuration)

// Operation is a function type representing an operation that may be retried.
//
// It encapsulates a task that returns an error to indicate success (nil) or failure (non-nil).
// This type is used with the Retry function for operations that do not produce a result.
//
// Returns:
//   - err (error): The error from the operation, or nil if the operation succeeded.
type Operation func() (err error)

// withEmptyData wraps an Operation to convert it into an OperationWithData that returns an
// empty struct as its data.
//
// This method enables operations that do not produce a result to be used with RetryWithData,
// allowing a consistent interface for both result-producing and non-result-producing operations.
//
// Returns:
//   - operationWithData (OperationWithData[struct{}]): A function that executes the original
//     Operation and returns an empty struct alongside the operation's error.
func (o Operation) withEmptyData() (operationWithData OperationWithData[struct{}]) {
	operationWithData = func() (struct{}, error) {
		return struct{}{}, o()
	}

	return
}

// OperationWithData is a generic function type representing an operation that returns both
// a result of type T and an error.
//
// It is used with RetryWithData to handle operations that produce a result, allowing the retrier
// to return the successful result alongside a nil error when the operation succeeds.
//
// Type Parameters:
//   - T: The type of the data returned by the operation.
//
// Returns:
//   - data (T): The result of the operation if successful.
//   - err (error): The error from the operation, or nil if the operation succeeded.
type OperationWithData[T any] func() (data T, err error)

// WithRetryMax returns an Option that sets the maximum number of retry attempts.
//
// It configures the retrier to limit retries to the specified number. Once this limit is reached,
// the retrier stops and returns the last error. A value of 0 means no retries are attempted
// (only the initial attempt is made).
//
// Parameters:
//   - retryMax (int): The maximum number of retry attempts. Should be non-negative; negative
//     values may lead to undefined behavior.
//
// Returns:
//   - option (Option): A functional option that sets the retryMax field in the configuration.
func WithRetryMax(retryMax int) (option Option) {
	return func(configuration *configuration) {
		configuration.retryMax = retryMax
	}
}

// WithRetryWaitMin returns an Option that sets the minimum delay between retry attempts.
//
// It defines the base delay for backoff calculations, ensuring retries do not occur too rapidly.
// This is particularly important for preventing overwhelming a system with rapid retries.
//
// Parameters:
//   - retryWaitMin (time.Duration): The minimum delay duration. Should be non-negative; negative
//     values may lead to undefined behavior.
//
// Returns:
//   - option (Option): A functional option that sets the retryWaitMin field in the configuration.
func WithRetryWaitMin(retryWaitMin time.Duration) (option Option) {
	return func(configuration *configuration) {
		configuration.retryWaitMin = retryWaitMin
	}
}

// WithRetryWaitMax returns an Option that sets the maximum delay between retry attempts.
//
// It caps the backoff duration to prevent excessively long delays, ensuring retries occur within
// a reasonable timeframe. Typically, retryWaitMax should be greater than or equal to retryWaitMin.
//
// Parameters:
//   - retryWaitMax (time.Duration): The maximum delay duration. Should be non-negative; negative
//     values may lead to undefined behavior.
//
// Returns:
//   - option (Option): A functional option that sets the retryWaitMax field in the configuration.
func WithRetryWaitMax(retryWaitMax time.Duration) (option Option) {
	return func(configuration *configuration) {
		configuration.retryWaitMax = retryWaitMax
	}
}

// WithRetryBackoff returns an Option that sets the backoff strategy for computing retry delays.
//
// It allows users to specify a custom backoff strategy (from the backoff package) to calculate
// delays based on the attempt number, minimum delay, and maximum delay. This enables flexible
// retry policies, such as exponential backoff with or without jitter.
//
// Parameters:
//   - retryBackoff (backoff.Backoff): The backoff strategy function. If nil, the retrier will
//     use a default strategy (e.g., exponential backoff).
//
// Returns:
//   - option (Option): A functional option that sets the retryBackoff field in the configuration.
func WithRetryBackoff(retryBackoff backoff.Backoff) (option Option) {
	return func(configuration *configuration) {
		configuration.retryBackoff = retryBackoff
	}
}

// WithNotifier returns an Option that sets a notifier callback for retry attempts.
//
// It configures a callback function that is invoked after each failed retry attempt, receiving
// the error and the computed backoff duration. This is useful for logging, monitoring, or other
// side effects during retries.
//
// Parameters:
//   - notifier (Notifier): The callback function to be called on each retry. If nil, no
//     notification is performed.
//
// Returns:
//   - option (Option): A functional option that sets the notifier field in the configuration.
func WithNotifier(notifier Notifier) (option Option) {
	return func(configuration *configuration) {
		configuration.notifier = notifier
	}
}

// Retry executes an operation with retries, respecting the provided context and configuration.
//
// It attempts the operation up to retryMax times (as specified in the configuration), waiting
// between attempts according to the backoff strategy. If the operation succeeds (returns nil
// error), it returns immediately. If the context is canceled or times out, it returns the
// context's error. If all retries fail, it returns the last error from the operation.
//
// Parameters:
//   - ctx (context.Context): The context controlling the retry lifecycle. Cancellation or timeout
//     aborts retries and returns ctx.Err().
//   - operation (Operation): The operation to retry, which returns an error indicating success
//     or failure.
//   - options (...Option): Variadic configuration options to customize retry behavior, such as
//     maximum retries, delay bounds, backoff strategy, and notifier.
//
// Returns:
//   - err (error): The error from the last attempt if all retries fail, or ctx.Err() if the
//     context is canceled or times out. Returns nil if the operation succeeds.
func Retry(ctx context.Context, operation Operation, options ...Option) (err error) {
	_, err = RetryWithData(ctx, operation.withEmptyData(), options...)

	return
}

// RetryWithData executes a generic operation that returns data and an error, with retries.
//
// It attempts the operation up to retryMax times, using the configured backoff strategy to
// compute delays between attempts. If the operation succeeds (returns nil error), it returns
// the operation's result and nil. If the context is canceled or times out, it returns the
// context's error. If all retries fail, it returns the last result and error from the operation.
//
// Parameters:
//   - ctx (context.Context): The context controlling the retry lifecycle. Cancellation or timeout
//     aborts retries and returns ctx.Err().
//   - operation (OperationWithData[T]): The operation to retry, returning a result of type T
//     and an error.
//   - options (...Option): Variadic configuration options to customize retry behavior.
//
// Returns:
//   - result (T): The result from the operation if it succeeds, or the last result if all retries fail.
//   - err (error): The error from the last attempt if all retries fail, or ctx.Err() if the
//     context is canceled or times out. Returns nil if the operation succeeds.
func RetryWithData[T any](ctx context.Context, operation OperationWithData[T], options ...Option) (result T, err error) {
	cfg := &configuration{
		retryMax:     3,
		retryWaitMin: 1 * time.Second,
		retryWaitMax: 30 * time.Second,
		retryBackoff: backoff.ExponentialWithDecorrelatedJitter(),
	}

	for _, option := range options {
		option(cfg)
	}

	for attempt := 1; ; attempt++ {
		select {
		case <-ctx.Done():
			err = ctx.Err()

			return
		default:
			result, err = operation()
			if err == nil {
				return
			}

			if cfg.retryMax > 0 && attempt >= cfg.retryMax {
				return
			}

			b := cfg.retryBackoff(cfg.retryWaitMin, cfg.retryWaitMax, attempt)

			if cfg.notifier != nil {
				cfg.notifier(err, b)
			}

			ticker := time.NewTicker(b)

			select {
			case <-ticker.C:
				ticker.Stop()
			case <-ctx.Done():
				ticker.Stop()

				err = context.Cause(ctx)

				return
			}
		}
	}
}
