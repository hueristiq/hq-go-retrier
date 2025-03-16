// Package backoff provides a collection of functions to generate backoff delays for retry strategies.
// These functions implement various exponential backoff algorithms — with and without jitter — to help
// manage retries in network operations, distributed systems, and other fault-tolerant applications.
//
// In many retry scenarios, it is beneficial to increase the delay between successive attempts in an
// exponential manner. However, using pure exponential backoff can lead to synchronization issues (i.e.
// the "thundering herd" problem) when multiple clients retry simultaneously. To mitigate this, jitter
// (i.e. randomness) is often added to the delay.
//
// Example Usage:
//
//	package main
//
//	import (
//	    "fmt"
//	    "time"
//	    "go.source.hueristiq.com/retrier/backoff"
//	)
//
//	func main() {
//	    backoffFunc := backoff.Exponential()
//	    delay := backoffFunc(1*time.Second, 30*time.Second, 3)
//	    fmt.Println("Retry delay:", delay)
//	}
package backoff
