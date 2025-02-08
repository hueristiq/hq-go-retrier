package retrier

import (
	"context"
	"time"

	"go.source.hueristiq.com/retrier/backoff"
)

// Operation is a function type that represents an operation that can be retried.
// The operation returns an error, which indicates whether the operation failed or succeeded.
//
// Returns:
//   - err (error): The error indicating whether the operation succeeded or failed.
type Operation func() (err error)

// withEmptyData wraps an Operation function to convert it into an OperationWithData that
// returns an empty struct. This is used for cases where the operation does not return any data
// but can be retried with the same mechanism as data-returning operations.
//
// Returns:
//   - operationWithData (OperationWithData[struct{}]): An OperationWithData function that returns an empty struct and error,
//     allowing non-data-returning operations to be handled by the RetryWithData function.
func (o Operation) withEmptyData() (operationWithData OperationWithData[struct{}]) {
	operationWithData = func() (struct{}, error) {
		return struct{}{}, o()
	}

	return
}

// OperationWithData is a function type that represents an operation that returns data along with an error.
// The generic type T allows the operation to return any type of data, making the retrier versatile for operations
// that may return results along with a possible error.
type OperationWithData[T any] func() (data T, err error)

// Retry attempts to execute the provided operation with a retry mechanism, using the provided options.
// If the operation continues to fail, it will retry based on the configuration, which may include max retries,
// backoff strategies, and min/max delay between retries.
//
// Arguments:
//   - ctx (context.Context): A context to control the lifetime of the retry operation. If the context is canceled or times out,
//     the retry operation will stop and return the context's error.
//   - operation (Operation): The operation to be retried.
//   - options (...Option): Optional configuration options that can adjust max retries, backoff strategy, or delay intervals.
//
// Returns:
//   - err (err): The error returned by the last failed attempt, or the context's error if the operation is canceled.
//
// Example:
//
//	err := retrier.Retry(ctx, someOperation, retrier.WithMaxRetries(5), retrier.WithBackoff(backoff.Exponential()))
//	// Retries 'someOperation' up to 5 times with exponential backoff.
func Retry(ctx context.Context, operation Operation, options ...Option) (err error) {
	_, err = RetryWithData(ctx, operation.withEmptyData(), options...)

	return
}

// RetryWithData attempts to execute the provided operation, which returns data along with an error, using the retry mechanism.
// It retries the operation based on the configuration and returns the result data if successful, or an error if the retries fail.
//
// Arguments:
//   - ctx (context.Context): A context to control the lifetime of the retry operation. If the context is canceled or times out,
//     the retry operation will stop and return the context's error.
//   - operation (OperationWithData[T]): The operation to be retried, which returns a value of type T and an error.
//   - options (...Option): Optional configuration options that can adjust max retries, backoff strategy, or delay intervals.
//
// Returns:
//   - result (T): The result of the operation if it succeeds within the allowed retry attempts.
//   - err (error): The error returned by the last failed attempt, or the context's error if the operation is canceled.
//
// Example:
//
//	result, err := retrier.RetryWithData(ctx, fetchData, retrier.WithMaxRetries(5), retrier.WithBackoff(backoff.Exponential()))
//	// Retries 'fetchData' up to 5 times with exponential backoff.
func RetryWithData[T any](ctx context.Context, operation OperationWithData[T], options ...Option) (result T, err error) {
	cfg := &Configuration{
		maxRetries: 3,
		maxDelay:   1000 * time.Millisecond,
		minDelay:   100 * time.Millisecond,
		backoff:    backoff.Exponential(),
	}

	for _, option := range options {
		option(cfg)
	}

	for attempt := range cfg.maxRetries {
		select {
		case <-ctx.Done():
			err = ctx.Err()

			return
		default:
			result, err = operation()
			if err == nil {
				return
			}

			b := cfg.backoff(cfg.minDelay, cfg.maxDelay, attempt)

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
