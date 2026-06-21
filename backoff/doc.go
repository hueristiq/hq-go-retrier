// Package backoff provides strategies for computing the delay between retry attempts.
//
// A strategy is any function matching the [Backoff] type. The package implements exponential
// backoff, where the base delay grows as minDelay * 2^attempt and is capped at maxDelay, in four
// variants that differ only in how they apply randomness (jitter) drawn from the
// [github.com/hueristiq/hq-lib-retrier-go/jitter] package:
//
//   - [Exponential] applies no jitter; the delay is fully deterministic.
//   - [ExponentialWithEqualJitter] keeps half the base fixed and randomizes the rest, yielding a
//     delay in [base/2, base].
//   - [ExponentialWithFullJitter] randomizes the entire base, yielding a delay in [0, base). This
//     spreads retries most aggressively.
//   - [ExponentialWithDecorrelatedJitter] draws from a range that widens with the previous delay,
//     decoupling successive waits.
//
// Jittered strategies reduce the "thundering herd" effect, where many clients that failed together
// retry in lockstep and overwhelm a recovering service. Each constructor returns a ready-to-use
// Backoff; pass one to [github.com/hueristiq/hq-lib-retrier-go.WithRetryBackoff].
//
// All strategies guard against integer overflow and return a zero duration for invalid input —
// non-positive bounds or a negative attempt.
package backoff
