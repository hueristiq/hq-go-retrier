// Package retrier runs operations that may fail transiently, retrying them with a configurable
// backoff until they succeed, a limit is reached, or the context is canceled.
//
// It suits work whose failures are often temporary — network requests, database queries, calls to
// external APIs — where a brief pause and another attempt is more useful than failing on the first
// error.
//
// # Entry points
//
// [Retry] runs an [Operation] that returns only an error. [RetryWithData] runs an
// [OperationWithData] that returns a value alongside its error and hands that value back to the
// caller. Retry is a thin wrapper over RetryWithData, so both share identical retry, backoff, and
// cancellation behavior.
//
// # Configuration
//
// Behavior is set through functional options passed to either entry point:
//
//   - [WithRetryMax] caps the number of attempts, counting the initial call. A value less than or
//     equal to zero retries indefinitely until success or cancellation.
//   - [WithRetryWaitMin] and [WithRetryWaitMax] bound the delay between attempts.
//   - [WithRetryBackoff] selects the strategy that computes each delay.
//   - [WithNotifier] registers a callback invoked after every failed attempt.
//
// Unset options fall back to defaults: three attempts, a one-second minimum and thirty-second
// maximum wait, and exponential backoff with decorrelated jitter.
//
// # Backoff and jitter
//
// The delay between attempts comes from a [github.com/hueristiq/hq-lib-retrier-go/backoff.Backoff]
// function. The backoff package provides exponential strategies, optionally combined with the
// jitter strategies in [github.com/hueristiq/hq-lib-retrier-go/jitter] to spread retries across
// clients and avoid the thundering-herd problem.
//
// # Context
//
// Every attempt and every wait observes the supplied context. When the context is canceled or its
// deadline passes, the active call returns [context.Cause] for that context and abandons any
// further attempts. On a clean context, the result of the final attempt — its value for
// RetryWithData, plus its error — is returned.
package retrier
