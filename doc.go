// Package retrier provides configurable retry logic for operations that may fail transiently.
// It encapsulates settings such as the number of retry attempts, delays between attempts,
// backoff strategies, and notification callbacks. This allows developers to easily customize
// the retry behavior for various operations, such as network requests, database transactions,
// or any operation that benefits from being retried upon failure.
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
//	    "go.source.hueristiq.com/retrier/backoff"
//	)
//
//	func main() {
//	    ctx := context.Background()
//	    err := retrier.Retry(ctx, someOperation,
//	        retrier.WithRetryMax(5),                                    // Allow a maximum of 5 retries.
//	        retrier.WithRetryWaitMin(100 * time.Millisecond),           // Set the minimum delay to 100ms.
//	        retrier.WithRetryWaitMax(2 * time.Second),                  // Set the maximum delay to 2 seconds.
//	        retrier.WithRetryBackoff(backoff.ExponentialWithFullJitter()), // Use exponential backoff with full jitter.
//	        retrier.WithNotifier(func(err error, backoff time.Duration) {
//	            // Log the error and the delay before the next retry.
//	            fmt.Printf("Retrying after error: %v, backoff: %v\n", err, backoff)
//	        }),
//	    )
//	    if err != nil {
//	        fmt.Println("Operation failed after retries:", err)
//	    }
//	}
//
// The retrier package offers a highly customizable and flexible approach to implementing retry logic.
// By leveraging configurable backoff strategies and notifier callbacks, you can tailor retry behavior
// to the specific requirements of your distributed systems or transient error-prone operations.
// This package abstracts the complexities of retry management, enabling developers to focus on
// the core logic of their applications.
package retrier
