# hq-lib-retrier-go

![made with go](https://img.shields.io/badge/made%20with-Go-1E90FF.svg) [![go reference](https://pkg.go.dev/badge/github.com/hueristiq/hq-lib-retrier-go.svg)](https://pkg.go.dev/github.com/hueristiq/hq-lib-retrier-go) [![go report card](https://goreportcard.com/badge/github.com/hueristiq/hq-lib-retrier-go)](https://goreportcard.com/report/github.com/hueristiq/hq-lib-retrier-go) [![license](https://img.shields.io/badge/license-MIT-gray.svg?color=1E90FF)](https://github.com/hueristiq/hq-lib-retrier-go/blob/master/LICENSE) ![maintenance](https://img.shields.io/badge/maintained%3F-yes-1E90FF.svg) [![open issues](https://img.shields.io/github/issues-raw/hueristiq/hq-lib-retrier-go.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-lib-retrier-go/issues?q=is:issue+is:open) [![closed issues](https://img.shields.io/github/issues-closed-raw/hueristiq/hq-lib-retrier-go.svg?style=flat&color=1E90FF)](https://github.com/hueristiq/hq-lib-retrier-go/issues?q=is:issue+is:closed) [![contribution](https://img.shields.io/badge/contributions-welcome-1E90FF.svg)](https://github.com/hueristiq/hq-lib-retrier-go/blob/master/CONTRIBUTING.md)

`hq-lib-retrier-go` is a [Go (Golang)](http://golang.org/) package for retrying operations that may fail transiently, such as network requests, database queries, or external API calls.

## Resource

- [Features](#features)
- [Installation](#installation)
- [Usage](#usage)
	- [Basic Retry](#basic-retry)
	- [Retry With Data](#retry-with-data)
- [Configuration](#configuration)
- [Backoff & Jitter Strategies](#backoff--jitter-strategies)
- [Contributing](#contributing)
- [Licensing](#licensing)

## Features

- **Configurable retry policy**: Set the maximum number of attempts, the minimum and maximum delay between attempts, and the backoff strategy.
- **Context support**: Every attempt and every wait observes the supplied `context.Context`, so cancellation and deadlines are respected immediately.
- **Result-carrying operations**: `RetryWithData` retries operations that return a value alongside an error and hands the value back to the caller.
- **Notifier callback**: A callback fires after each failed attempt with the triggering error and the next delay — useful for logging, metrics, or debugging.
- **Backoff and jitter strategies**: Built-in exponential backoff with equal, full, or decorrelated jitter to mitigate the "thundering herd" problem in distributed systems.

## Installation

To install `hq-lib-retrier-go`, run the following command in your Go project:

```bash
go get -v -u github.com/hueristiq/hq-lib-retrier-go
```

This package requires Go 1.26.3 or later.

## Usage

The package exposes two entry points: `Retry` for operations that return only an error, and `RetryWithData` for operations that return both a value and an error. Both accept the same functional options.

### Basic Retry

Use `Retry` for operations that return only an error, such as network requests or file operations:

```go
package main

import (
	"context"
	"fmt"
	"time"

	hqgoretrier "github.com/hueristiq/hq-lib-retrier-go"
	"github.com/hueristiq/hq-lib-retrier-go/backoff"
)

func main() {
	attempts := 0

	operation := func() error {
		attempts++

		if attempts < 3 {
			return fmt.Errorf("temporary failure (attempt %d)", attempts)
		}

		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err := hqgoretrier.Retry(ctx, operation,
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryWaitMin(100*time.Millisecond),
		hqgoretrier.WithRetryWaitMax(2*time.Second),
		hqgoretrier.WithRetryBackoff(backoff.ExponentialWithFullJitter()),
		hqgoretrier.WithNotifier(func(err error, next time.Duration) {
			fmt.Printf("Retry due to error: %v. Next attempt in %v.\n", err, next)
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

Use `RetryWithData` for operations that return both a result and an error, such as fetching data from an API. The result type is inferred from the operation, so no type assertions are needed:

```go
package main

import (
	"context"
	"fmt"
	"time"

	hqgoretrier "github.com/hueristiq/hq-lib-retrier-go"
	"github.com/hueristiq/hq-lib-retrier-go/backoff"
)

func main() {
	attempts := 0

	fetchData := func() (string, error) {
		attempts++

		if attempts < 2 {
			return "", fmt.Errorf("failed to fetch data (attempt %d)", attempts)
		}

		return "payload", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	result, err := hqgoretrier.RetryWithData(ctx, fetchData,
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryWaitMin(200*time.Millisecond),
		hqgoretrier.WithRetryWaitMax(3*time.Second),
		hqgoretrier.WithRetryBackoff(backoff.Exponential()),
		hqgoretrier.WithNotifier(func(err error, next time.Duration) {
			fmt.Printf("Retrying after error: %v, waiting: %v\n", err, next)
		}),
	)
	if err != nil {
		fmt.Printf("Failed to fetch data after retries: %v\n", err)

		return
	}

	fmt.Printf("Data fetched successfully: %s\n", result)
}
```

When the context is canceled or its deadline passes, both functions stop retrying and return [`context.Cause(ctx)`](https://pkg.go.dev/context#Cause). On a clean context, the result of the final attempt — its value for `RetryWithData`, plus its error — is returned.

## Configuration

Behavior is set through functional options passed to `Retry` or `RetryWithData`. Any option left unset falls back to its default.

| Option | Description | Default |
| --- | --- | --- |
| `WithRetryMax(n)` | Maximum number of attempts, **including the initial call** (so `3` means the first call plus up to two retries). A value `<= 0` retries indefinitely until success or context cancellation. | `3` |
| `WithRetryWaitMin(d)` | Lower bound for the delay between attempts. A value `<= 0` falls back to the default. | `1s` |
| `WithRetryWaitMax(d)` | Upper bound for the delay between attempts. A value `<= 0` falls back to the default; a value below the minimum is raised to the minimum. | `30s` |
| `WithRetryBackoff(b)` | Strategy that computes each delay (see below). Passing `nil` selects the default. | `ExponentialWithDecorrelatedJitter()` |
| `WithNotifier(fn)` | Callback invoked after each failed attempt with the error and the next delay. | none |

Invalid values are normalized before the first attempt, so a misconfigured retrier never spins in a zero-delay loop.

## Backoff & Jitter Strategies

The delay between attempts is produced by a `backoff.Backoff` function. The `backoff` package provides exponential strategies, where the base delay grows as `min(maxDelay, minDelay * 2^attempt)`. Jitter adds randomness so that clients which failed together do not retry in lockstep.

| Strategy | Delay range | Notes |
| --- | --- | --- |
| `Exponential()` | `base` | Deterministic; no jitter. |
| `ExponentialWithEqualJitter()` | `[base/2, base)` | Half the delay is fixed, half is random. |
| `ExponentialWithFullJitter()` | `[0, base)` | Fully randomized; spreads retries most aggressively. |
| `ExponentialWithDecorrelatedJitter()` | `[minDelay, min(maxDelay, previous*3)]` | Default; decouples successive delays. Stateful. |

Here `base` is `min(maxDelay, minDelay * 2^attempt)` and `previous` is the delay the strategy produced on its preceding call — the decorrelated strategy remembers its last delay, so each draw depends on the previous random draw (true decorrelated jitter). It is safe for concurrent use, but a shared instance couples the delay sequences of all its callers; create one per retry loop for independent decorrelation. All strategies guard against integer overflow and return a zero duration for invalid input (non-positive bounds, a minimum above the maximum, or a negative attempt).

The jitter functions are also exported directly from the [`jitter`](https://pkg.go.dev/github.com/hueristiq/hq-lib-retrier-go/jitter) package for building custom strategies. To implement your own, supply any function matching the `backoff.Backoff` signature to `WithRetryBackoff`:

```go
hqgoretrier.WithRetryBackoff(func(minDelay, maxDelay time.Duration, attempt int) time.Duration {
	return minDelay // constant backoff
})
```

## Contributing

Contributions are welcome and encouraged! Feel free to submit [Pull Requests](https://github.com/hueristiq/hq-lib-retrier-go/pulls) or report [Issues](https://github.com/hueristiq/hq-lib-retrier-go/issues). For more details, check out the [contribution guidelines](https://github.com/hueristiq/hq-lib-retrier-go/blob/master/CONTRIBUTING.md).

A big thank you to all the [contributors](https://github.com/hueristiq/hq-lib-retrier-go/graphs/contributors) for your ongoing support!

![contributors](https://contrib.rocks/image?repo=hueristiq/hq-lib-retrier-go&max=500)

## Licensing

This package is licensed under the [MIT license](https://opensource.org/license/mit). You are free to use, modify, and distribute it, as long as you follow the terms of the license. You can find the full license text in the repository - [Full MIT license text](https://github.com/hueristiq/hq-lib-retrier-go/blob/master/LICENSE).
