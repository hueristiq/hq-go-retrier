// Package backoff provides strategies for implementing retry mechanisms with controlled delays.
//
// This package includes various exponential backoff strategies, incorporating jitter techniques to
// prevent synchronized retries, reduce congestion, and improve system resilience.
//
// Available Strategies:
//   - Exponential: Basic exponential backoff with a growth factor.
//   - ExponentialWithEqualJitter: Exponential backoff with equal jitter for moderate randomness.
//   - ExponentialWithFullJitter: Exponential backoff with full jitter for maximum randomness.
//   - ExponentialWithDecorrelatedJitter: Exponential backoff with decorrelated jitter to prevent
//     uncontrolled exponential growth.
//
// These strategies are useful for retrying failed operations in distributed systems, API calls,
// and network requests where implementing controlled delays enhances system stability.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "time"
//	    "go.source.hueristiq.com/retrier/backoff"
//	)
//
//	func main() {
//	    backoffFunc := backoff.Exponential()
//	    delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	    fmt.Println("Retry delay:", delay)
//	}
//
// Reference:
// - https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
//
// This package is designed to optimize retry logic by introducing controlled delays with jitter,
// reducing congestion and improving overall system efficiency.
package backoff
