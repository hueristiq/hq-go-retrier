// Package jitter provides functions for applying random jitter to backoff durations.
// Jitter is a technique used to add randomness to retry intervals, preventing
// thundering herd issues and reducing load on resources in distributed systems.
//
// In distributed systems or network environments, multiple clients may retry
// failed operations simultaneously, leading to increased load on the system.
// By introducing jitter, the retry intervals are randomized, reducing the
// chance of synchronized retries (the "thundering herd" problem) and
// distributing the system load more evenly over time.
//
// This package offers three jitter strategies:
//  1. **Equal Jitter**: Adds moderate randomness to the retry interval by
//     selecting a random value within half of the original backoff duration.
//     This is suitable for scenarios where you want some consistency with
//     slight randomization.
//  2. **Full Jitter**: Introduces maximum randomness by selecting a retry
//     interval uniformly between 0 and the original backoff duration. This
//     is ideal for highly distributed systems where maximum variation is
//     preferred.
//  3. **Decorrelated Jitter**: Produces a random backoff duration influenced
//     by the previous backoff value, keeping the retry interval bounded
//     within a specified range. This is useful for preventing unbounded
//     exponential growth in retry delays.
package jitter
