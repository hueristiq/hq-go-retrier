package retrier

import (
	"context"
	"time"

	"github.com/hueristiq/hq-go-retrier/backoff"
)

// Operation is a function type that represents an operation that can be retried.
// The operation returns an error, which indicates whether the operation failed or succeeded.
type Operation func() (err error)

// withEmptyData wraps an Operation function to convert it into an OperationWithData that
// returns an empty struct. This is used for cases where the operation does not return any data
// but can be retried with the same mechanism as data-returning operations.
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
// Parameters:
//   - ctx: A context to control the lifetime of the retry operation. If the context is canceled or times out,
//     the retry operation will stop and return the context's error.
//   - operation: The operation to be retried.
//   - opts: Optional configuration options that can adjust max retries, backoff strategy, or delay intervals.
//
// Returns:
//   - err: The error returned by the last failed attempt, or the context's error if the operation is canceled.
//
// Example:
//
//	err := retrier.Retry(ctx, someOperation, retrier.WithMaxRetries(5), retrier.WithBackoff(backoff.Exponential()))
//	// Retries 'someOperation' up to 5 times with exponential backoff.
func Retry(ctx context.Context, operation Operation, opts ...Option) (err error) {
	_, err = RetryWithData(ctx, operation.withEmptyData(), opts...)

	return
}

// RetryWithData attempts to execute the provided operation, which returns data along with an error, using the retry mechanism.
// It retries the operation based on the configuration and returns the result data if successful, or an error if the retries fail.
//
// Parameters:
//   - ctx: A context to control the lifetime of the retry operation. If the context is canceled or times out,
//     the retry operation will stop and return the context's error.
//   - operation: The operation to be retried, which returns a value of type T and an error.
//   - opts: Optional configuration options that can adjust max retries, backoff strategy, or delay intervals.
//
// Returns:
//   - result: The result of the operation if it succeeds within the allowed retry attempts.
//   - err: The error returned by the last failed attempt, or the context's error if the operation is canceled.
//
// Example:
//
//	result, err := retrier.RetryWithData(ctx, fetchData, retrier.WithMaxRetries(5), retrier.WithBackoff(backoff.Exponential()))
//	// Retries 'fetchData' up to 5 times with exponential backoff.
func RetryWithData[T any](ctx context.Context, operation OperationWithData[T], opts ...Option) (result T, err error) {
	// Set default retry configuration.
	cfg := &Configuration{
		maxRetries: 3,                       // Default maximum retries
		maxDelay:   1000 * time.Millisecond, // Default maximum delay between retries
		minDelay:   100 * time.Millisecond,  // Default minimum delay between retries
		backoff:    backoff.Exponential(),   // Default backoff strategy: exponential
	}

	// Apply any provided options to configure retry behavior.
	for _, opt := range opts {
		opt(cfg)
	}

	// Retry loop up to maxRetries.
	for attempt := range cfg.maxRetries {
		select {
		case <-ctx.Done():
			// If the context is done, return the context's error.
			err = ctx.Err()

			return
		default:
			// Execute the operation and check for success.
			result, err = operation()
			if err == nil {
				// Operation succeeded, return the result.
				return
			}

			// If the operation fails, calculate the backoff delay.
			b := cfg.backoff(cfg.minDelay, cfg.maxDelay, attempt)

			// Wait for the backoff period before the next retry attempt.
			ticker := time.NewTicker(b)

			select {
			case <-ticker.C:
				// Backoff delay is over, stop the ticker and retry.
				ticker.Stop()
			case <-ctx.Done():
				// Context is done, stop the ticker and return the context's error.
				ticker.Stop()

				err = ctx.Err()

				return
			}
		}
	}

	return
}
