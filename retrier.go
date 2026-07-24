package retrier

import (
	"context"
	"time"

	"github.com/hueristiq/hq-lib-retrier-go/backoff"
)

const (
	defaultRetryMax     = 3
	defaultRetryWaitMin = 1 * time.Second
	defaultRetryWaitMax = 30 * time.Second
)

// options holds the settings for retry operations, defining the behavior of the retry
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
type options struct {
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

// OptionFunc is a function type used to modify the retry options in a declarative manner.
//
// It allows users to customize retry behavior by setting fields in the options struct,
// such as the maximum number of retries, delay bounds, backoff strategy, or notifier callback.
// Multiple options can be combined to create a tailored retry policy.
//
// Parameters:
//   - opts (*options): A pointer to the options struct to be modified.
type OptionFunc func(opts *options)

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

// WithRetryMax returns an OptionFunc that sets the maximum number of attempts.
//
// The count includes the initial attempt, so a value of 3 means one initial call followed by up to
// two retries. Once the limit is reached, the retrier stops and returns the last error. A value
// less than or equal to 0 means the operation is retried indefinitely until it succeeds or the
// context is canceled.
//
// Parameters:
//   - retryMax (int): The maximum number of attempts, including the initial one. Values less than
//     or equal to 0 enable unlimited retries.
//
// Returns:
//   - (OptionFunc): A functional option that sets the retryMax field in the options.
func WithRetryMax(retryMax int) OptionFunc {
	return func(opts *options) {
		opts.retryMax = retryMax
	}
}

// WithRetryWaitMin returns an OptionFunc that sets the minimum delay between retry attempts.
//
// It defines the base delay for backoff calculations, ensuring retries do not occur too rapidly.
// This is particularly important for preventing overwhelming a system with rapid retries.
//
// Parameters:
//   - retryWaitMin (time.Duration): The minimum delay duration. Values less than or equal to 0
//     fall back to the default of 1 second, so a misconfigured retrier cannot spin in a
//     zero-delay loop.
//
// Returns:
//   - (OptionFunc): A functional option that sets the retryWaitMin field in the options.
func WithRetryWaitMin(retryWaitMin time.Duration) OptionFunc {
	return func(opts *options) {
		opts.retryWaitMin = retryWaitMin
	}
}

// WithRetryWaitMax returns an OptionFunc that sets the maximum delay between retry attempts.
//
// It caps the backoff duration to prevent excessively long delays, ensuring retries occur within
// a reasonable timeframe. Typically, retryWaitMax should be greater than or equal to retryWaitMin.
//
// Parameters:
//   - retryWaitMax (time.Duration): The maximum delay duration. Values less than or equal to 0
//     fall back to the default of 30 seconds, and a value below retryWaitMin is raised to
//     retryWaitMin.
//
// Returns:
//   - (OptionFunc): A functional option that sets the retryWaitMax field in the options.
func WithRetryWaitMax(retryWaitMax time.Duration) OptionFunc {
	return func(opts *options) {
		opts.retryWaitMax = retryWaitMax
	}
}

// WithRetryBackoff returns an OptionFunc that sets the backoff strategy for computing retry delays.
//
// It allows users to specify a custom backoff strategy (from the backoff package) to calculate
// delays based on the attempt number, minimum delay, and maximum delay. This enables flexible
// retry policies, such as exponential backoff with or without jitter.
//
// Parameters:
//   - retryBackoff (backoff.Backoff): The backoff strategy function. If nil, the default
//     strategy (exponential backoff with decorrelated jitter) is used.
//
// Returns:
//   - (OptionFunc): A functional option that sets the retryBackoff field in the options.
func WithRetryBackoff(retryBackoff backoff.Backoff) OptionFunc {
	return func(opts *options) {
		opts.retryBackoff = retryBackoff
	}
}

// WithNotifier returns an OptionFunc that sets a notifier callback for retry attempts.
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
//   - (OptionFunc): A functional option that sets the notifier field in the options.
func WithNotifier(notifier Notifier) OptionFunc {
	return func(opts *options) {
		opts.notifier = notifier
	}
}

// Retry executes an operation with retries, respecting the provided context and options.
//
// It attempts the operation up to retryMax times (as specified in the options), waiting
// between attempts according to the backoff strategy. If the operation succeeds (returns nil
// error), it returns immediately. If the context is canceled or times out, it returns the
// context's error. If all retries fail, it returns the last error from the operation.
//
// Parameters:
//   - ctx (context.Context): The context controlling the retry lifecycle. Cancellation or timeout
//     aborts retries and returns ctx.Err().
//   - operation (Operation): The operation to retry, which returns an error indicating success
//     or failure.
//   - ofs (...OptionFunc): Variadic functional options to customize retry behavior, such as
//     maximum attempts, delay bounds, backoff strategy, and notifier.
//
// Returns:
//   - err (error): The error from the last attempt if all retries fail, or ctx.Err() if the
//     context is canceled or times out. Returns nil if the operation succeeds.
func Retry(ctx context.Context, operation Operation, ofs ...OptionFunc) (err error) {
	_, err = RetryWithData(ctx, operation.withEmptyData(), ofs...)

	return
}

// RetryWithData executes a generic operation that returns data and an error, with retries.
//
// It attempts the operation up to retryMax times, using the configured backoff strategy to
// compute delays between attempts. If the operation succeeds (returns nil error), it returns
// the operation's result and nil. If the context is canceled or times out, it returns the
// context's error. If all retries fail, it returns the last result and error from the operation.
//
// After the options are applied, invalid values are normalized: a nil backoff falls back to
// exponential backoff with decorrelated jitter, non-positive wait bounds fall back to their
// defaults, and a retryWaitMax below retryWaitMin is raised to retryWaitMin. This guarantees
// the retry loop never spins with a zero delay.
//
// Parameters:
//   - ctx (context.Context): The context controlling the retry lifecycle. Cancellation or timeout
//     aborts retries and returns ctx.Err().
//   - operation (OperationWithData[T]): The operation to retry, returning a result of type T
//     and an error.
//   - ofs (...OptionFunc): Variadic functional options to customize retry behavior.
//
// Returns:
//   - result (T): The result from the operation if it succeeds, or the last result if all retries fail.
//   - err (error): The error from the last attempt if all retries fail, or ctx.Err() if the
//     context is canceled or times out. Returns nil if the operation succeeds.
func RetryWithData[T any](ctx context.Context, operation OperationWithData[T], ofs ...OptionFunc) (result T, err error) {
	opts := &options{
		retryMax:     defaultRetryMax,
		retryWaitMin: defaultRetryWaitMin,
		retryWaitMax: defaultRetryWaitMax,
		retryBackoff: backoff.ExponentialWithDecorrelatedJitter(),
	}

	for _, f := range ofs {
		f(opts)
	}

	if opts.retryBackoff == nil {
		opts.retryBackoff = backoff.ExponentialWithDecorrelatedJitter()
	}

	if opts.retryWaitMin <= 0 {
		opts.retryWaitMin = defaultRetryWaitMin
	}

	if opts.retryWaitMax <= 0 {
		opts.retryWaitMax = defaultRetryWaitMax
	}

	if opts.retryWaitMax < opts.retryWaitMin {
		opts.retryWaitMax = opts.retryWaitMin
	}

	for attempt := 1; ; attempt++ {
		if ctx.Err() != nil {
			err = context.Cause(ctx)

			return
		}

		result, err = operation()
		if err == nil {
			return
		}

		if opts.retryMax > 0 && attempt >= opts.retryMax {
			return
		}

		b := opts.retryBackoff(opts.retryWaitMin, opts.retryWaitMax, attempt)

		if opts.notifier != nil {
			opts.notifier(err, b)
		}

		timer := time.NewTimer(b)

		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()

			err = context.Cause(ctx)

			return
		}
	}
}
