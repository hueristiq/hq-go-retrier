// Package jitter provides various strategies for introducing randomized delays in retry mechanisms.
//
// This package implements multiple jitter strategies, including:
// - Equal Jitter: Balances between fixed backoff and randomization.
// - Full Jitter: Provides a completely randomized delay within the range.
// - Decorrelated Jitter: Prevents uncontrolled exponential growth while maintaining randomness.
//
// These strategies help in mitigating synchronized retry bursts, which can cause excessive load on systems.
// The jitter functions can be useful in distributed systems, network communication, and backoff policies.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "time"
//	    "go.source.hueristiq.com/retrier/jitter"
//	)
//
//	func main() {
//	    backoff := 10 * time.Second
//	    jitteredBackoff := jitter.Equal(backoff)
//	    fmt.Println("Jittered Backoff:", jitteredBackoff)
//	}
//
// Features:
// - Provides structured jitter strategies for retry mechanisms.
// - Implements cryptographic randomness for secure random backoff durations.
// - Ensures retry intervals are varied to reduce request collisions.
//
// Reference:
// - https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
//
// This package is designed to improve resilience and efficiency in retry mechanisms.
package jitter
