// Package retrier provides a flexible mechanism for retrying operations with configurable
// backoff strategies and delay settings. It allows users to specify how many times an operation
// should be retried, the delay between retries, and the strategy used for increasing the delay
// (e.g., exponential backoff with jitter).
//
// This package is useful in scenarios where transient failures may occur, such as network
// requests, database operations, or other I/O-related tasks. The configurable retry mechanism
// helps in ensuring resilience by retrying operations with controlled delays, preventing overwhelming
// the system with constant retries.
package retrier
