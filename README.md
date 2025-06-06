# hq-go-retrier

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/hq-go-retrier)](https://goreportcard.com/report/github.com/hueristiq/hq-go-retrier) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-retrier.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-retrier.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md)

`hq-go-retrier` is a [Go (Golang)](http://golang.org/) package for managing retries for operations that might temporarily fail, such as network requests, database queries, or external API calls e.t.c.

## Resource

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
	- [Basic Retry](#basic-retry)
	- [Retry With Data](#retry-with-data)
- [Contributing](#contributing)
- [Licensing](#licensing)

## Features

- **Configurable Retry Mechanism**: Supports number of retries, minimum and maximum delay durations, and backoff strategies customization to suit your application's needs.
- **Context Support**: Integrates with Go's `context` package for graceful cancellation and timeout management, ensuring retries respect the application's lifecycle constraints.
- **Data Handling**: Supports operations that return data (via `RetryWithData`), enabling retries for operations with results, not just errors.
- **Notifier Callback**: ProSupports defination of a notifier callback mechanism to monitor retry attempts, useful for logging, metrics, or debugging.
- **Backoff and Jitter Strategies**: Includes built-in backoff and jitter strategies to mitigate the "thundering herd" problem in distributed systems.

## Installation

To install `hq-go-retrier`, run the following command in your Go project:

```bash
go get -v -u github.com/hueristiq/hq-go-retrier
```

Make sure your Go environment is set up properly (Go 1.18 or later is recommended).

## Usage

The `hq-go-retrier` package provides two primary functions: `Retry` for operations that return only an error, and `RetryWithData` for operations that return both data and an error. Below are examples demonstrating various use cases.

### Basic Retry

Use the `Retry` function for operations that return only an error, such as network requests or file operations:

```go
package main

import (
	"context"
	"fmt"
	"time"

	hqgoretrier "github.com/hueristiq/hq-go-retrier"
	"github.com/hueristiq/hq-go-retrier/backoff"
)

func main() {
	operation := func() error {
		return fmt.Errorf("temporary failure")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := hqgoretrier.Retry(ctx, operation,
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryWaitMin(100*time.Millisecond),
		hqgoretrier.WithRetryWaitMax(2*time.Second),
		hqgoretrier.WithRetryBackoff(backoff.ExponentialWithFullJitter()),
		hqgoretrier.WithNotifier(func(err error, b time.Duration) {
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

Use `RetryWithData` for operations that return both a result and an error, such as fetching data from an API:

```go
package main

import (
	"context"
	"fmt"
	"time"

	hqgoretrier "github.com/hueristiq/hq-go-retrier"
	"github.com/hueristiq/hq-go-retrier/backoff"
)

func fetchData() (string, error) {
	return "", fmt.Errorf("failed to fetch data")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	result, err := hqgoretrier.RetryWithData(ctx, fetchData,
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryWaitMin(200*time.Millisecond),
		hqgoretrier.WithRetryWaitMax(3*time.Second),
		hqgoretrier.WithRetryBackoff(backoff.Exponential()),
		hqgoretrier.WithNotifier(func(err error, b time.Duration) {
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

Contributions are welcome and encouraged! Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-retrier/pulls) or report [Issues](https://github.com/hueristiq/hq-go-retrier/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md).

A big thank you to all the [contributors](https://github.com/hueristiq/hq-go-retrier/graphs/contributors) for your ongoing support!

![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-retrier&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE).