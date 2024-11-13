# hq-go-retrier

![made with go](https://img.shields.io/badge/made%20with-Go-0000FF.svg) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=0000FF)](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0000FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-go-retrier.svg?style=flat&color=0000FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-go-retrier.svg?style=flat&color=0000FF)](https://github.com/hueristiq/hq-go-retrier/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-0000FF.svg)](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md)

`hq-go-retrier` is a [Go (Golang)](http://golang.org/) package designed to manage retries for operations that might temporarily fail. It allows developers to customize how retries are handled using different strategies, such as increasing the wait time between each attempt - backoffs and jitters.

> [!TIP]
> **Backoff** is a strategy used to manage retry intervals when handling transient failures in a system. Instead of retrying an operation immediately after a failure, the backoff mechanism increases the waiting period between retries, often to prevent overloading the system or further exacerbating the issue.
>
> **Jitter** is a technique used in conjunction with backoff strategies to introduce randomness to the retry intervals. Its main goal is to avoid the "thundering herd" problem, where multiple clients or processes attempt to retry a failed operation at the same time, overwhelming the system or service they're interacting with.

## Resource

* [Features](#features)
* [Installation](#installation)
* [Usage](#usage)
	* [Configuration Options](#configuration-options)
* [Contributing](#contributing)
* [Licensing](#licensing)
* [Credits](#credits)
	* [Contributors](#contributors)
	* [Similar Projects](#similar-projects)

## Features

* **Configurable Retry Mechanism:** Easily configure the maximum number of retries, minimum and maximum delays, and backoff strategies.
* **Custom Backoff Strategies:** Supports various backoff strategies, including exponential backoff and jitter to manage retries effectively.
* **Context Support:** Operations can be run with a context to handle cancellation and timeouts gracefully.
* **Data Handling:** Supports operations that return both data and error, enhancing its usability.

## Installation

To install the package, run the following command in your terminal:

```bash
go get -v -u github.com/hueristiq/hq-go-retrier
```

This command will download and install the `hq-go-retrier` package into your Go workspace, making it available for use in your projects.

## Usage

Here's a simple example demonstrating how to use `hq-go-retrier`:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	retrier "github.com/hueristiq/hq-go-retrier"
	"github.com/hueristiq/hq-go-retrier/backoff"
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

### Configuration Options

The following options can be used to customize the retry behavior:

* `WithMaxRetries(int)`: Sets the maximum number of retry attempts.
* `WithMinDelay(time.Duration)`: Sets the minimum delay between retries.
* `WithMaxDelay(time.Duration)`: Sets the maximum delay between retries.
* `WithBackoff(backoff.Backoff)`: Sets the backoff strategy to be used.
* `WithNotifier(notifier)`: Sets a callback function that gets triggered on each retry attempt, providing feedback on errors and backoff.

## Contributing

We welcome contributions! Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-go-retrier/pulls) or report [Issues](https://github.com/hueristiq/hq-go-retrier/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-go-retrier/blob/master/CONTRIBUTING.md).

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-go-retrier/blob/master/LICENSE).

## Credits

### Contributors

A huge thanks to all the contributors who have helped make `hq-go-retrier` what it is today!

[![contributors](https://contrib.rocks/image?repo=hueristiq/hq-go-retrier&max=500)](https://github.com/hueristiq/hq-go-retrier/graphs/contributors)

### Similar Projects

If you're interested in more packages like this, check out:

[Cenk Alti's backoff](https://github.com/cenkalti/backoff) ◇ [Avast's retry-go](https://github.com/avast/retry-go)