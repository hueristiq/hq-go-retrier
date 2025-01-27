package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.source.hueristiq.com/retrier"
	"go.source.hueristiq.com/retrier/backoff"
)

func main() {
	operation := func() error {
		// Simulate a failing operation
		fmt.Println("Trying operation...")
		return errors.New("operation failed")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Retry the operation with custom configuration
	err := retrier.Retry(ctx, operation,
		retrier.WithMaxRetries(5),
		retrier.WithMinDelay(100*time.Millisecond),
		retrier.WithMaxDelay(1*time.Second),
		retrier.WithBackoff(backoff.ExponentialWithDecorrelatedJitter()),
		retrier.WithNotifier(func(err error, backoff time.Duration) {
			fmt.Printf("Operation failed: %v\n", err)
			fmt.Printf("...wait %d seconds for the next retry\n\n", backoff)
		}),
	)

	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	} else {
		fmt.Println("Operation succeeded")
	}
}
