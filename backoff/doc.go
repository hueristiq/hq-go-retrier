// Package backoff provides a collection of functions to generate backoff delays for retry strategies.
// These functions implement various exponential backoff algorithms — with and without jitter — to help
// manage retries in network operations, distributed systems, and other fault-tolerant applications.
//
// In many retry scenarios, it is beneficial to increase the delay between successive attempts in an
// exponential manner. However, using pure exponential backoff can lead to synchronization issues (i.e.
// the "thundering herd" problem) when multiple clients retry simultaneously. To mitigate this, jitter
// (i.e. randomness) is often added to the delay.
//
// This package defines several functions that return a backoff function with the following signature:
//
// func(minDelay, maxDelay time.Duration, attempt int) (delay time.Duration)
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
