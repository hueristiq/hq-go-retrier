// Package retrier provides a flexible and configurable retry mechanism for transient failures.
//
// This package allows developers to implement structured retry logic with customizable parameters,
// such as maximum retries, minimum and maximum delay durations, and backoff strategies. It also
// includes a notifier callback function for tracking retry attempts.
//
// The retrier is useful in scenarios where operations may intermittently fail, such as network requests,
// database queries, or distributed system interactions. By applying controlled backoff and retry logic,
// applications can enhance their resilience and prevent unnecessary resource exhaustion.
//
// Features:
// - Configurable retry policies with min/max delay and max retry count.
// - Supports multiple backoff strategies, including exponential backoff with jitter.
// - Allows the use of a notifier callback function for logging retry attempts.
// - Context-aware execution to allow cancellation of retry operations.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "context"
//	    "fmt"
//	    "time"
//	    "go.source.hueristiq.com/retrier"
//	    "go.source.hueristiq.com/retrier/retrier"
//	    "go.source.hueristiq.com/retrier/backoff"
//	)
//
//	func main() {
//	    ctx := context.Background()
//	    err := retrier.Retry(ctx, someOperation,
//	        retrier.WithMaxRetries(5),
//	        retrier.WithBackoff(backoff.ExponentialWithFullJitter()),
//	        retrier.WithNotifier(func(err error, backoff time.Duration) {
//	            fmt.Printf("Retrying after error: %v, backoff: %v\n", err, backoff)
//	        }),
//	    )
//	    if err != nil {
//	        fmt.Println("Operation failed after retries:", err)
//	    }
//	}
//
// Reference:
// - https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
//
// This package is designed to improve reliability in distributed systems by implementing structured
// retry mechanisms that prevent overwhelming dependent services and optimize request handling.
package retrier
