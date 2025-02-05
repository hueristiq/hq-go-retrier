# hq-go-retrier

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/xsubfind3r)](https://goreportcard.com/report/github.com/hueristiq/hq-go-retrier) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-retrier.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-retrier.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md)

`hq-go-retrier` is a [Go (Golang)](http://golang.org/) package for managing retries for operations that might temporarily fail. It allows developers to customize how retries are handled using different strategies.

## Resource

* [Features](#features)
* [Usage](#usage)
* [Contributing](#contributing)
* [Licensing](#licensing)

## Features

* **Configurable Retry Mechanism:** Easily configure the maximum number of retries, minimum and maximum delays, and backoff strategies.
* **Custom Backoff Strategies:** Supports various backoff strategies, including exponential backoff and jitter to manage retries effectively.
* **Context Support:** Operations can be run with a context to handle cancellation and timeouts gracefully.
* **Data Handling:** Supports operations that return both data and error, enhancing its usability.

## Usage

```bash
go get -v -u go.source.hueristiq.com/retrier
```

Here's a simple example demonstrating how to use `hq-go-retrier`:

```go
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
```

The following options can be used to customize the retry behavior:

* `WithMaxRetries(int)`: Sets the maximum number of retry attempts.
* `WithMinDelay(time.Duration)`: Sets the minimum delay between retries.
* `WithMaxDelay(time.Duration)`: Sets the maximum delay between retries.
* `WithBackoff(backoff.Backoff)`: Sets the backoff strategy to be used.
* `WithNotifier(notifier)`: Sets a callback function that gets triggered on each retry attempt, providing feedback on errors and backoff.

## Contributing

Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-retrier/pulls) or report [Issues](https://github.com/hueristiq/hq-go-retrier/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md).

Huge thanks to the [contributors](https://github.com/hueristiq/hq-go-retrier/graphs/contributors) thus far!

![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-retrier&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE).