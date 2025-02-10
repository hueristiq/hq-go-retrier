package retrier

import (
	"context"
	"time"

	"go.source.hueristiq.com/retrier/backoff"
)

// Operation is a function type that represents an operation which may be retried.
// The function returns an error which indicates whether the operation succeeded (nil error)
// or failed (non-nil error).
//
// Returns:
//   - err (error): The error encountered during the operation, or nil if successful.
type Operation func() (err error)

// withEmptyData wraps an Operation to convert it into an OperationWithData that returns an
// empty struct as its data. This is useful for retrying operations that do not produce a
// result, allowing them to be used with RetryWithData.
//
// Returns:
//   - operationWithData (OperationWithData[struct{}]): A function that, when invoked,
//     returns an empty struct and the error from the original Operation.
func (o Operation) withEmptyData() (operationWithData OperationWithData[struct{}]) {
	operationWithData = func() (struct{}, error) {
		return struct{}{}, o()
	}

	return
}

// OperationWithData is a generic function type representing an operation that returns data
// of type T along with an error. This allows the retrier to handle operations that produce a result.
type OperationWithData[T any] func() (data T, err error)

// Retry attempts to execute the provided operation (of type Operation) with a retry mechanism.
// It uses the supplied context to control the lifetime of the retry loop. If the operation fails,
// Retry will continue attempting until the operation succeeds, the maximum number of retries is
// reached, or the context is canceled.
//
// Arguments:
//   - ctx (context.Context): The context controlling the retry operation. If canceled or timed out,
//     Retry will stop and return ctx.Err().
//   - operation (Operation): The operation to be retried. A nil error indicates success.
//   - options (...Option): A variadic list of configuration options to customize the retry behavior
//     (e.g., maximum retries, delay intervals, backoff strategy, notifier).
//
// Returns:
//   - err (error): The error from the last retry attempt if the operation never succeeded, or the
//     context's error if the context is canceled.
//
// Example:
//
//	err := retrier.Retry(ctx, someOperation, retrier.WithRetryMax(5), retrier.WithRetryBackoff(backoff.Exponential()))
//	if err != nil {
//	    // handle error
//	}
func Retry(ctx context.Context, operation Operation, options ...Option) (err error) {
	_, err = RetryWithData(ctx, operation.withEmptyData(), options...)

	return
}

// RetryWithData attempts to execute the provided operation (which returns data and an error)
// using a retry mechanism. It applies the configured retry settings and returns the result of the
// operation if it eventually succeeds. If all retry attempts fail, it returns the error from the
// last attempt (or the context's error if the context is canceled).
//
// Arguments:
//   - ctx (context.Context): The context that controls the retry lifecycle. Cancellation or timeout
//     will abort further retries.
//   - operation (OperationWithData[T]): The operation to be retried, returning a value of type T and an error.
//   - options (...Option): Optional configuration options to customize the retry policy.
//
// Returns:
//   - result (T): The result of the operation if successful.
//   - err (error): The error from the final attempt if the operation never succeeds, or the context's error
//     if the context is canceled.
//
// Example:
//
//	result, err := retrier.RetryWithData(ctx, fetchData, retrier.WithRetryMax(5), retrier.WithRetryBackoff(backoff.Exponential()))
//	if err != nil {
//	    // handle error
//	} else {
//	    // process result
//	}
func RetryWithData[T any](ctx context.Context, operation OperationWithData[T], options ...Option) (result T, err error) {
	cfg := &Configuration{
		retryMax:     3,
		retryWaitMin: 1 * time.Second,
		retryWaitMax: 30 * time.Second,
		retryBackoff: backoff.Exponential(),
	}

	for _, option := range options {
		option(cfg)
	}

	for attempt := range cfg.retryMax {
		select {
		case <-ctx.Done():
			err = ctx.Err()

			return
		default:
			result, err = operation()
			if err == nil {
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

				err = ctx.Err()

				return
			}
		}
	}

	return
}
