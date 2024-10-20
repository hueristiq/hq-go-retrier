// Package backoff provides various strategies for calculating retry backoff intervals
// with optional jitter. Backoff strategies are commonly used in systems where
// retrying failed operations is necessary, such as network requests, database
// connections, or API calls. By progressively increasing the wait time between
// retries, backoff mechanisms can reduce load on resources and prevent overwhelming
// a system that may be experiencing temporary failure.
//
// This package supports multiple exponential backoff strategies, including:
//  1. **Exponential Backoff**: A strategy where the delay between retries increases
//     exponentially based on the number of attempts.
//  2. **Exponential Backoff with Equal Jitter**: Adds a moderate amount of randomness
//     to the exponential delay, making the retry interval less predictable.
//  3. **Exponential Backoff with Full Jitter**: Applies maximum randomness to the
//     retry interval, introducing full jitter to the exponential delay.
//  4. **Exponential Backoff with Decorrelated Jitter**: Calculates the retry interval
//     based on the previous delay, ensuring bounded and random backoff durations.
//
// By adding jitter, the retry intervals are randomized, preventing the "thundering herd"
// problem where multiple clients retry operations simultaneously, leading to further
// system overload.
package backoff
