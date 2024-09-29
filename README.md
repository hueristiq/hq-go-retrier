# hqgoretry

![made with go](https://img.shields.io/badge/made%20with-Go-0000FF.svg) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=0000FF)](https://github.com/hueristiq/hqgoretry/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-0000FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hqgoretry.svg?style=flat&color=0000FF)](https://github.com/hueristiq/hqgoretry/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hqgoretry.svg?style=flat&color=0000FF)](https://github.com/hueristiq/hqgoretry/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-0000FF.svg)](https://github.com/hueristiq/hqgoretry/blob/master/CONTRIBUTING.md)

`hqgoretry` is a lightweight and flexible [Go (Golang)](http://golang.org/) package designed to implement retry mechanisms with customizable backoff strategies. It helps developers manage transient failures efficiently by retrying operations that may fail intermittently.

> [!TIP]
> **Backoff** is a strategy used to manage retry intervals when handling transient failures in a system. Instead of retrying an operation immediately after a failure, the backoff mechanism increases the waiting period between retries, often to prevent overloading the system or further exacerbating the issue.
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
go get -v -u github.com/hueristiq/hqgoretry
```

This command will download and install the `hqgoretry` package into your Go workspace, making it available for use in your projects.

## Usage

Here's a simple example demonstrating how to use `hqgoretry`:

```go
package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hueristiq/hqgoretry"
	"github.com/hueristiq/hqgoretry/backoff"
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
	err := hqgoretry.Retry(ctx, operation,
		hqgoretry.WithMaxRetries(5),
		hqgoretry.WithMinDelay(100*time.Millisecond),
		hqgoretry.WithMaxDelay(1*time.Second),
		hqgoretry.WithBackoff(backoff.ExponentialWithDecorrelatedJitter()),
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

## Contributing

We welcome contributions! Feel free to submit [Pull Requests](https://github.com/hueristiq/hqgoretry/pulls) or report [Issues](https://github.com/hueristiq/hqgoretry/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hqgoretry/blob/master/CONTRIBUTING.md).

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hqgoretry/blob/master/LICENSE).

## Credits

### Contributors

A huge thanks to all the contributors who have helped make `hqgoretry` what it is today!

[![contributors](https://contrib.rocks/image?repo=hueristiq/hqgoretry&max=500)](https://github.com/hueristiq/hqgoretry/graphs/contributors)

### Similar Projects

If you're interested in more packages like this, check out:

[Cenk Alti's backoff](https://github.com/cenkalti/backoff) ◇ [Avast's retry-go](https://github.com/avast/retry-go)