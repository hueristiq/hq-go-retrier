# hq-go-retrier

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/xsubfind3r)](https://goreportcard.com/report/github.com/hueristiq/hq-go-retrier) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-retrier.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-retrier.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md)

`hq-go-retrier` is a [Go (Golang)](http://golang.org/) package for managing retries for operations that might temporarily fail. It allows you to easily implement retry logic with customizable parameters such as the number of retries, delays between retries, backoff strategies, and notifications on retry attempts. This package is especially useful in networked, distributed, or fault-tolerant applications where transient errors are common.

## Resource

* [Features](#features)
* [Usage](#usage)
	* [Basic Retry](#basic-retry)
	* [Retry With Data](#retry-with-data)
* [Contributing](#contributing)
* [Licensing](#licensing)

## Features

* **Configurable Retry Mechanism:** Easily configure the maximum number of retries, minimum and maximum delays, and backoff strategies.
* **Flexible Backoff Strategies:** Supports various backoff strategies, including exponential backoff and jitter to manage retries effectively.
* **Notifier Callback:** Optionally receive notifications for each retry attempt with details about the error and the backoff duration.
* **Context Support:** Operations can be run with a context to handle cancellation and timeouts gracefully.
* **Data Handling:** Supports operations that return both data and error, enhancing its usability.

## Usage

```bash
go get -v -u go.source.hueristiq.com/retrier
```

### Basic Retry

The simplest usage of `hq-go-retrier` is to retry an operation that only returns an error. Use the `Retry` function along with any optional configuration options:

```go
package main

import (
	"context"
	"fmt"
	"time"

	"go.source.hueristiq.com/retrier"
	"go.source.hueristiq.com/retrier/backoff"
)

func main() {
	// Define an operation that may fail.
	operation := func() error {
		// Replace with your logic that might fail.
		return fmt.Errorf("an error occurred")
	}

	// Create a context with timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retry the operation with custom configuration.
	err := retrier.Retry(ctx, operation,
		retrier.WithRetryMax(5),                                    // Maximum 5 retries.
		retrier.WithRetryWaitMin(100*time.Millisecond),             // Minimum wait of 100ms.
		retrier.WithRetryWaitMax(2*time.Second),                    // Maximum wait of 2 seconds.
		retrier.WithRetryBackoff(backoff.ExponentialWithFullJitter()),// Exponential backoff with full jitter.
		retrier.WithNotifier(func(err error, b time.Duration) {
			fmt.Printf("Retry due to error: %v. Next attempt in %v.\n", err, b)
		}),
	)
	if err != nil {
		fmt.Printf("Operation failed after retries: %v\n", err)
	} else {
		fmt.Println("Operation succeeded!")
	}
}
```

### Retry With Data

If the operation returns data along with an error, use `RetryWithData`. This function allows to obtain the result from the operation once it succeeds.

```go
package main

import (
	"context"
	"fmt"
	"time"

	"go.source.hueristiq.com/retrier"
	"go.source.hueristiq.com/retrier/backoff"
)

// fetchData simulates an operation that returns a string result.
func fetchData() (string, error) {
	// Replace with your logic. For example:
	return "", fmt.Errorf("failed to fetch data")
}

func main() {
	// Create a context.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Retry the operation that returns data.
	result, err := retrier.RetryWithData(ctx, fetchData,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(200*time.Millisecond),
		retrier.WithRetryWaitMax(3*time.Second),
		retrier.WithRetryBackoff(backoff.Exponential()),
		retrier.WithNotifier(func(err error, b time.Duration) {
			fmt.Printf("Retrying after error: %v, waiting: %v\n", err, b)
		}),
	)
	if err != nil {
		fmt.Printf("Failed to fetch data after retries: %v\n", err)
		return
	}
	fmt.Printf("Data fetched successfully: %s\n", result)
}
```

## Contributing

Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-retrier/pulls) or report [Issues](https://github.com/hueristiq/hq-go-retrier/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md).

Huge thanks to the [contributors](https://github.com/hueristiq/hq-go-retrier/graphs/contributors) thus far!

![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-retrier&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE).