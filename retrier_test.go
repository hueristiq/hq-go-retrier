package retrier_test

import (
	"context"
	"errors"
	"math"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hqgoretrier "github.com/hueristiq/hq-lib-retrier-go"
	hqgoretrierbackoff "github.com/hueristiq/hq-lib-retrier-go/backoff"
)

var (
	errTestOperation = errors.New("operation failed")
	errAlternate     = errors.New("alternate failure")
)

type flakyOperation struct {
	failures int
	calls    int
}

func (o *flakyOperation) run() error {
	o.calls++

	if o.calls <= o.failures {
		return errTestOperation
	}

	return nil
}

func noWaitBackoff(_, _ time.Duration) hqgoretrierbackoff.Backoff {
	return func(int) time.Duration {
		return 0
	}
}

func TestRetry_ImmediateSuccess(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 0}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.NoError(t, err, "Expected the operation to succeed on the first attempt")
	assert.Equal(t, 1, op.calls, "Expected the operation to be called exactly once")
}

func TestRetry_SuccessAfterFailures(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 2}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.NoError(t, err, "Expected the operation to succeed after retries")
	assert.Equal(t, 3, op.calls, "Expected the operation to be called 3 times")
}

func TestRetry_AcrossBackoffStrategies(t *testing.T) {
	t.Parallel()

	strategies := []struct {
		name       string
		newBackoff func(minDelay, maxDelay time.Duration) hqgoretrierbackoff.Backoff
	}{
		{"exponential", hqgoretrierbackoff.Exponential},
		{"equal jitter", hqgoretrierbackoff.ExponentialWithEqualJitter},
		{"full jitter", hqgoretrierbackoff.ExponentialWithFullJitter},
		{"decorrelated jitter", hqgoretrierbackoff.ExponentialWithDecorrelatedJitter},
	}

	for _, s := range strategies {
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			op := &flakyOperation{failures: 2}

			err := hqgoretrier.Retry(
				t.Context(),
				op.run,
				hqgoretrier.WithMaxAttempts(5),
				hqgoretrier.WithRetryWaitMin(time.Millisecond),
				hqgoretrier.WithRetryWaitMax(5*time.Millisecond),
				hqgoretrier.WithRetryBackoff(s.newBackoff),
			)

			require.NoError(t, err, "Expected the operation to succeed after retries")
			assert.Equal(t, 3, op.calls, "Expected the operation to be called 3 times")
		})
	}
}

func TestRetry_MaxAttempts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		maxAttempts int
		wantCalls   int
	}{
		{"single attempt", 1, 1},
		{"two attempts", 2, 2},
		{"three attempts", 3, 3},
		{"five attempts", 5, 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			op := &flakyOperation{failures: 1000}

			err := hqgoretrier.Retry(
				t.Context(),
				op.run,
				hqgoretrier.WithMaxAttempts(tt.maxAttempts),
				hqgoretrier.WithRetryBackoff(noWaitBackoff),
			)

			require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
			assert.Equal(t, tt.wantCalls, op.calls, "Expected attempts to equal maxAttempts")
		})
	}
}

func TestRetry_MaxAttemptsFallsBackToDefault(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		maxAttempts int
	}{
		{"zero", 0},
		{"negative", -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			op := &flakyOperation{failures: 1000}

			err := hqgoretrier.Retry(
				t.Context(),
				op.run,
				hqgoretrier.WithMaxAttempts(tt.maxAttempts),
				hqgoretrier.WithRetryBackoff(noWaitBackoff),
			)

			require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
			assert.Equal(t, 3, op.calls, "Expected a non-positive maxAttempts to fall back to the default of 3 attempts")
		})
	}
}

func TestRetry_RetryMaxDeprecatedAlias(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 1000}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithRetryMax(2),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
	assert.Equal(t, 2, op.calls, "Expected the deprecated alias to forward to WithMaxAttempts")
}

func TestRetry_OptionsLastOneWins(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 1000}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithMaxAttempts(2),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
	assert.Equal(t, 2, op.calls, "Expected the later option to override the earlier one")
}

func TestRetry_DefaultMaxAttempts(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 1000}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.Error(t, err, "Expected the operation to fail after the default number of attempts")
	assert.Equal(t, 3, op.calls, "Expected the default of 3 attempts")
}

func TestRetry_ReturnsLastError(t *testing.T) {
	t.Parallel()

	calls := 0

	op := func() error {
		calls++

		if calls == 1 {
			return errTestOperation
		}

		return errAlternate
	}

	err := hqgoretrier.Retry(
		t.Context(),
		op,
		hqgoretrier.WithMaxAttempts(2),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, errAlternate, "Expected the error from the final attempt")
	assert.NotErrorIs(t, err, errTestOperation, "Expected the earlier error to be discarded")
}

func TestRetry_RetryIf(t *testing.T) {
	t.Parallel()

	t.Run("rejected error stops retries", func(t *testing.T) {
		t.Parallel()

		op := &flakyOperation{failures: 1000}

		err := hqgoretrier.Retry(
			t.Context(),
			op.run,
			hqgoretrier.WithMaxAttempts(5),
			hqgoretrier.WithRetryBackoff(noWaitBackoff),
			hqgoretrier.WithRetryIf(func(err error) bool {
				return errors.Is(err, errAlternate)
			}),
		)

		require.ErrorIs(t, err, errTestOperation, "Expected the rejected error")
		assert.Equal(t, 1, op.calls, "Expected no retries once the predicate rejects the error")
	})

	t.Run("accepted error is retried", func(t *testing.T) {
		t.Parallel()

		op := &flakyOperation{failures: 2}

		err := hqgoretrier.Retry(
			t.Context(),
			op.run,
			hqgoretrier.WithMaxAttempts(5),
			hqgoretrier.WithRetryBackoff(noWaitBackoff),
			hqgoretrier.WithRetryIf(func(err error) bool {
				return errors.Is(err, errTestOperation)
			}),
		)

		require.NoError(t, err, "Expected the operation to succeed after retries")
		assert.Equal(t, 3, op.calls, "Expected the operation to be called 3 times")
	})

	t.Run("nil predicate retries every error", func(t *testing.T) {
		t.Parallel()

		op := &flakyOperation{failures: 2}

		err := hqgoretrier.Retry(
			t.Context(),
			op.run,
			hqgoretrier.WithMaxAttempts(5),
			hqgoretrier.WithRetryBackoff(noWaitBackoff),
			hqgoretrier.WithRetryIf(nil),
		)

		require.NoError(t, err, "Expected the operation to succeed after retries")
		assert.Equal(t, 3, op.calls, "Expected the operation to be called 3 times")
	})
}

func TestRetry_NotifierReceivesAttemptErrorAndWait(t *testing.T) {
	t.Parallel()

	type notification struct {
		attempt int
		err     error
		wait    time.Duration
	}

	var got []notification

	perAttempt := func(_, _ time.Duration) hqgoretrierbackoff.Backoff {
		return func(attempt int) time.Duration {
			return time.Duration(attempt) * time.Millisecond
		}
	}

	op := &flakyOperation{failures: 2}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(perAttempt),
		hqgoretrier.WithNotifier(func(attempt int, err error, wait time.Duration) {
			got = append(got, notification{attempt: attempt, err: err, wait: wait})
		}),
	)

	require.NoError(t, err, "Expected the operation to succeed after retries")
	require.Len(t, got, 2, "Expected the notifier to fire once per failed attempt that is retried")

	for i, n := range got {
		assert.Equal(t, i+1, n.attempt, "Expected the 1-based number of the failed attempt")
		require.ErrorIs(t, n.err, errTestOperation, "Expected the notifier to receive the failure")
		assert.Equal(t, time.Duration(i+1)*time.Millisecond, n.wait, "Expected the computed delay")
	}
}

func TestRetry_NotifierSkipsFinalAttempt(t *testing.T) {
	t.Parallel()

	calls := 0

	op := &flakyOperation{failures: 1000}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithMaxAttempts(3),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
		hqgoretrier.WithNotifier(func(int, error, time.Duration) {
			calls++
		}),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
	assert.Equal(t, 2, calls, "Expected the notifier to fire only before the waits, not after the final failure")
}

func TestRetry_NotifierSilentWithoutUpcomingRetry(t *testing.T) {
	t.Parallel()

	t.Run("permanent error", func(t *testing.T) {
		t.Parallel()

		notified := 0

		op := func() error {
			return hqgoretrier.Permanent(errTestOperation)
		}

		err := hqgoretrier.Retry(
			t.Context(),
			op,
			hqgoretrier.WithMaxAttempts(5),
			hqgoretrier.WithRetryBackoff(noWaitBackoff),
			hqgoretrier.WithNotifier(func(int, error, time.Duration) {
				notified++
			}),
		)

		require.ErrorIs(t, err, errTestOperation, "Expected the permanent error")
		assert.Zero(t, notified, "Expected no notification when no retry follows a permanent error")
	})

	t.Run("rejected by predicate", func(t *testing.T) {
		t.Parallel()

		notified := 0

		op := &flakyOperation{failures: 1000}

		err := hqgoretrier.Retry(
			t.Context(),
			op.run,
			hqgoretrier.WithMaxAttempts(5),
			hqgoretrier.WithRetryBackoff(noWaitBackoff),
			hqgoretrier.WithRetryIf(func(error) bool { return false }),
			hqgoretrier.WithNotifier(func(int, error, time.Duration) {
				notified++
			}),
		)

		require.ErrorIs(t, err, errTestOperation, "Expected the rejected error")
		assert.Zero(t, notified, "Expected no notification when no retry follows a rejected error")
	})
}

func TestRetry_BackoffConstructedOnceWithConfiguredBounds(t *testing.T) {
	t.Parallel()

	var (
		constructs     int
		gotMin, gotMax time.Duration
		attempts       []int
	)

	waitMin := 7 * time.Millisecond
	waitMax := 11 * time.Millisecond

	factory := func(minDelay, maxDelay time.Duration) hqgoretrierbackoff.Backoff {
		constructs++
		gotMin, gotMax = minDelay, maxDelay

		return func(attempt int) time.Duration {
			attempts = append(attempts, attempt)

			return 0
		}
	}

	op := &flakyOperation{failures: 3}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryWaitMin(waitMin),
		hqgoretrier.WithRetryWaitMax(waitMax),
		hqgoretrier.WithRetryBackoff(factory),
	)

	require.NoError(t, err, "Expected the operation to succeed after retries")
	assert.Equal(t, 1, constructs, "Expected the backoff to be constructed once per retry loop")
	assert.Equal(t, waitMin, gotMin, "Expected the configured minimum wait")
	assert.Equal(t, waitMax, gotMax, "Expected the configured maximum wait")
	assert.Equal(t, []int{1, 2, 3}, attempts, "Expected increasing attempt numbers")
}

func TestRetry_NilBackoffFallsBackToDefault(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)

	defer cancel()

	op := &flakyOperation{failures: 1000}

	err := hqgoretrier.Retry(
		ctx,
		op.run,
		hqgoretrier.WithRetryBackoff(nil),
	)

	require.ErrorIs(t, err, context.DeadlineExceeded, "Expected the default backoff's wait to outlast the context")
	assert.Equal(t, 1, op.calls, "Expected a single attempt before the default backoff's wait begins")
}

func TestRetry_WaitBoundsAreNormalized(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		waitMin time.Duration
		waitMax time.Duration
		wantMin time.Duration
		wantMax time.Duration
	}{
		{"zero bounds", 0, 0, time.Second, 30 * time.Second},
		{"negative bounds", -time.Second, -time.Second, time.Second, 30 * time.Second},
		{"zero min", 0, 5 * time.Second, time.Second, 5 * time.Second},
		{"zero max", 5 * time.Second, 0, 5 * time.Second, 30 * time.Second},
		{"max below min", 10 * time.Second, 5 * time.Second, 10 * time.Second, 10 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var gotMin, gotMax time.Duration

			factory := func(minDelay, maxDelay time.Duration) hqgoretrierbackoff.Backoff {
				gotMin, gotMax = minDelay, maxDelay

				return func(int) time.Duration {
					return 0
				}
			}

			op := &flakyOperation{failures: 1000}

			err := hqgoretrier.Retry(
				t.Context(),
				op.run,
				hqgoretrier.WithMaxAttempts(2),
				hqgoretrier.WithRetryWaitMin(tt.waitMin),
				hqgoretrier.WithRetryWaitMax(tt.waitMax),
				hqgoretrier.WithRetryBackoff(factory),
			)

			require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
			assert.Equal(t, tt.wantMin, gotMin, "Expected the normalized minimum wait")
			assert.Equal(t, tt.wantMax, gotMax, "Expected the normalized maximum wait")
		})
	}
}

func TestRetry_ContextCanceledBeforeStart(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	op := &flakyOperation{failures: 2}

	err := hqgoretrier.Retry(
		ctx,
		op.run,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, context.Canceled, "Expected a cancellation error")
	assert.Zero(t, op.calls, "Expected the operation never to run once the context is already canceled")
}

func TestRetry_ContextCanceledBeforeStartWithCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("shutting down")

	ctx, cancel := context.WithCancelCause(t.Context())
	cancel(cause)

	op := &flakyOperation{failures: 2}

	err := hqgoretrier.Retry(
		ctx,
		op.run,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, cause, "Expected the cancellation cause from the pre-attempt branch")
	assert.Zero(t, op.calls, "Expected the operation never to run once the context is already canceled")
}

func TestRetry_ContextCanceledDuringWait(t *testing.T) {
	t.Parallel()

	cause := errors.New("canceled during wait")

	ctx, cancel := context.WithCancelCause(t.Context())
	defer cancel(nil)

	waiting := make(chan struct{})

	var once sync.Once

	longBackoff := func(_, _ time.Duration) hqgoretrierbackoff.Backoff {
		return func(int) time.Duration {
			once.Do(func() { close(waiting) })

			return time.Hour
		}
	}

	op := func() error {
		return errTestOperation
	}

	go func() {
		<-waiting

		cancel(cause)
	}()

	err := hqgoretrier.Retry(
		ctx,
		op,
		hqgoretrier.WithMaxAttempts(10),
		hqgoretrier.WithRetryBackoff(longBackoff),
	)

	require.ErrorIs(t, err, cause, "Expected the cancellation cause from the wait branch")
}

func TestRetry_ContextTimeout(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 1000}

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)

	defer cancel()

	err := hqgoretrier.Retry(
		ctx,
		op.run,
		hqgoretrier.WithMaxAttempts(1000),
		hqgoretrier.WithRetryWaitMin(30*time.Millisecond),
		hqgoretrier.WithRetryWaitMax(100*time.Millisecond),
		hqgoretrier.WithRetryBackoff(hqgoretrierbackoff.Exponential),
	)

	require.ErrorIs(t, err, context.DeadlineExceeded, "Expected a deadline-exceeded error")
	assert.LessOrEqual(t, op.calls, 2, "Expected few attempts before the deadline elapses")
}

func TestRetry_UnboundedAttemptsStopOnContext(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 1_000_000}

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)

	defer cancel()

	err := hqgoretrier.Retry(
		ctx,
		op.run,
		hqgoretrier.WithMaxAttempts(math.MaxInt),
		hqgoretrier.WithRetryWaitMin(time.Millisecond),
		hqgoretrier.WithRetryWaitMax(2*time.Millisecond),
		hqgoretrier.WithRetryBackoff(hqgoretrierbackoff.Exponential),
	)

	require.ErrorIs(t, err, context.DeadlineExceeded, "Expected the context deadline to stop the retries")
	assert.Positive(t, op.calls, "Expected the operation to be attempted at least once")
}

func TestRetry_ConcurrentLoops(t *testing.T) {
	t.Parallel()

	const loops = 8

	var (
		wg       sync.WaitGroup
		failures atomic.Int32
	)

	for range loops {
		wg.Add(1)

		go func() {
			defer wg.Done()

			op := &flakyOperation{failures: 2}

			err := hqgoretrier.Retry(
				t.Context(),
				op.run,
				hqgoretrier.WithMaxAttempts(5),
				hqgoretrier.WithRetryWaitMin(time.Nanosecond),
				hqgoretrier.WithRetryWaitMax(time.Microsecond),
				hqgoretrier.WithRetryBackoff(hqgoretrierbackoff.ExponentialWithDecorrelatedJitter),
			)
			if err != nil {
				failures.Add(1)
			}
		}()
	}

	wg.Wait()

	assert.Zero(t, failures.Load(), "Expected all concurrent retry loops to succeed with per-loop backoff state")
}

func TestRetryWithData_ImmediateSuccess(t *testing.T) {
	t.Parallel()

	calls := 0

	op := func() (string, error) {
		calls++

		return "ok", nil
	}

	result, err := hqgoretrier.RetryWithData(
		t.Context(),
		op,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.NoError(t, err, "Expected the operation to succeed immediately")
	assert.Equal(t, "ok", result, "Expected the successful result")
	assert.Equal(t, 1, calls, "Expected the operation to be called exactly once")
}

func TestRetryWithData_SuccessAfterFailures(t *testing.T) {
	t.Parallel()

	calls := 0

	op := func() (int, error) {
		calls++

		if calls <= 2 {
			return 0, errTestOperation
		}

		return 42, nil
	}

	result, err := hqgoretrier.RetryWithData(
		t.Context(),
		op,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.NoError(t, err, "Expected the operation to succeed after retries")
	assert.Equal(t, 42, result, "Expected the successful result")
	assert.Equal(t, 3, calls, "Expected the operation to be called 3 times")
}

func TestRetryWithData_ReturnsLastResultOnFailure(t *testing.T) {
	t.Parallel()

	calls := 0

	op := func() (int, error) {
		calls++

		return calls, errTestOperation
	}

	result, err := hqgoretrier.RetryWithData(
		t.Context(),
		op,
		hqgoretrier.WithMaxAttempts(3),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
	assert.Equal(t, 3, result, "Expected the data from the final attempt")
	assert.Equal(t, 3, calls, "Expected three attempts")
}

func TestRetryWithData_RetryIfStopsRetries(t *testing.T) {
	t.Parallel()

	calls := 0

	op := func() (int, error) {
		calls++

		return calls, errTestOperation
	}

	result, err := hqgoretrier.RetryWithData(
		t.Context(),
		op,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
		hqgoretrier.WithRetryIf(func(err error) bool {
			return errors.Is(err, errAlternate)
		}),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the rejected error")
	assert.Equal(t, 1, result, "Expected the data from the only attempt")
	assert.Equal(t, 1, calls, "Expected no retries once the predicate rejects the error")
}

func TestRetryWithData_ContextCanceledBeforeStart(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	calls := 0

	op := func() (int, error) {
		calls++

		return 99, errTestOperation
	}

	result, err := hqgoretrier.RetryWithData(
		ctx,
		op,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, context.Canceled, "Expected a cancellation error")
	assert.Zero(t, result, "Expected the zero value when canceled before any attempt")
	assert.Zero(t, calls, "Expected the operation never to run")
}

func TestRetryWithData_ContextCanceledDuringWait(t *testing.T) {
	t.Parallel()

	cause := errors.New("canceled during wait")

	ctx, cancel := context.WithCancelCause(t.Context())
	defer cancel(nil)

	waiting := make(chan struct{})

	var once sync.Once

	longBackoff := func(_, _ time.Duration) hqgoretrierbackoff.Backoff {
		return func(int) time.Duration {
			once.Do(func() { close(waiting) })

			return time.Hour
		}
	}

	op := func() (int, error) {
		return 0, errTestOperation
	}

	go func() {
		<-waiting

		cancel(cause)
	}()

	result, err := hqgoretrier.RetryWithData(
		ctx,
		op,
		hqgoretrier.WithMaxAttempts(10),
		hqgoretrier.WithRetryBackoff(longBackoff),
	)

	require.ErrorIs(t, err, cause, "Expected the cancellation cause from the wait branch")
	assert.Zero(t, result, "Expected the zero value when canceled during the wait")
}
