// Package jitter provides various strategies for applying randomness (jitter) to backoff durations.
// Jitter is commonly used in retry algorithms to prevent synchronization between multiple clients
// (reducing the "thundering herd" problem) and to introduce variation in retry timing.
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
package jitter
