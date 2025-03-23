# hq-go-retrier

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/xsubfind3r)](https://goreportcard.com/report/github.com/hueristiq/hq-go-retrier) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-retrier.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-retrier.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md)

`hq-go-retrier` is a [Go (Golang)](http://golang.org/) package for managing retries for operations that might temporarily fail. It allows you to easily implement retry logic with customizable parameters such as the number of retries, delays between retries, backoff strategies, and notifications on retry attempts.

## Resource

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
	- [Basic Retry](#basic-retry)
	- [Retry With Data](#retry-with-data)
- [Contributing](#contributing)
- [Licensing](#licensing)

## Features

- **Configurable Retry Mechanism:** Customize the maximum number of retries, as well as the delay durations between retries.
- **Flexible Backoff Strategies:** Choose from various backoff and jitter strategies.
- **Context Support:** Integrates with Go's context package to support cancellation and timeouts.
- **Data Handling:** Use retry functions for both simple error-handling and operations that return valuable data.
- **Notifier Callback:** Get real-time notifications on each retry attempt to log errors or trigger custom actions.

## Installation

To install `hq-go-retrier`, run:

```bash
go get -v -u go.source.hueristiq.com/retrier
```

Make sure your Go environment is set up properly (Go 1.x or later is recommended).

## Usage

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
	operation := func() error {
		return fmt.Errorf("an error occurred")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := retrier.Retry(ctx, operation,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(100*time.Millisecond),
		retrier.WithRetryWaitMax(2*time.Second),
		retrier.WithRetryBackoff(backoff.ExponentialWithFullJitter()),
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

func fetchData() (string, error) {
	return "", fmt.Errorf("failed to fetch data")
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

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

Contributions are welcome and encouraged! Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-retrier/pulls) or report [Issues](https://github.com/hueristiq/hq-go-retrier/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md).

A big thank you to all the [contributors](https://github.com/hueristiq/hq-go-retrier/graphs/contributors) for your support!

![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-retrier&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE).