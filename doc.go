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
//	    hqgoretrier "github.com/hueristiq/hq-go-retrier"
//	    "github.com/hueristiq/hq-go-retrier/backoff"
//	)
//
//	func main() {
//	    ctx := context.Background()
//	    err := hqgoretrier.Retry(ctx, someOperation,
//	        hqgoretrier.WithRetryMax(5),                                    // Allow a maximum of 5 retries.
//	        hqgoretrier.WithRetryWaitMin(100 * time.Millisecond),           // Set the minimum delay to 100ms.
//	        hqgoretrier.WithRetryWaitMax(2 * time.Second),                  // Set the maximum delay to 2 seconds.
//	        hqgoretrier.WithRetryBackoff(backoff.ExponentialWithFullJitter()), // Use exponential backoff with full jitter.
//	        hqgoretrier.WithNotifier(func(err error, backoff time.Duration) {
//	            // Log the error and the delay before the next retry.
//	            fmt.Printf("Retrying after error: %v, backoff: %v\n", err, backoff)
//	        }),
//	    )
//	    if err != nil {
//	        fmt.Println("Operation failed after retries:", err)
//	    }
//	}
package retrier
