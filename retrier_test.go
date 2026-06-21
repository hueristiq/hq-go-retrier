package retrier_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	hqgoretrier "github.com/hueristiq/hq-lib-retrier-go"
	"github.com/hueristiq/hq-lib-retrier-go/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func noWaitBackoff() backoff.Backoff {
	return func(_, _ time.Duration, _ int) time.Duration {
		return 0
	}
}

func TestRetry_ImmediateSuccess(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 0}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
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
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
	)

	require.NoError(t, err, "Expected the operation to succeed after retries")
	assert.Equal(t, 3, op.calls, "Expected the operation to be called 3 times")
}

func TestRetry_AcrossBackoffStrategies(t *testing.T) {
	t.Parallel()

	strategies := []struct {
		name string
		b    backoff.Backoff
	}{
		{"exponential", backoff.Exponential()},
		{"equal jitter", backoff.ExponentialWithEqualJitter()},
		{"full jitter", backoff.ExponentialWithFullJitter()},
		{"decorrelated jitter", backoff.ExponentialWithDecorrelatedJitter()},
	}

	for _, s := range strategies {
		t.Run(s.name, func(t *testing.T) {
			t.Parallel()

			op := &flakyOperation{failures: 2}

			err := hqgoretrier.Retry(
				t.Context(),
				op.run,
				hqgoretrier.WithRetryMax(5),
				hqgoretrier.WithRetryWaitMin(time.Millisecond),
				hqgoretrier.WithRetryWaitMax(5*time.Millisecond),
				hqgoretrier.WithRetryBackoff(s.b),
			)

			require.NoError(t, err, "Expected the operation to succeed after retries")
			assert.Equal(t, 3, op.calls, "Expected the operation to be called 3 times")
		})
	}
}

func TestRetry_MaxAttempts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		retryMax  int
		wantCalls int
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
				hqgoretrier.WithRetryMax(tt.retryMax),
				hqgoretrier.WithRetryBackoff(noWaitBackoff()),
			)

			require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
			assert.Equal(t, tt.wantCalls, op.calls, "Expected attempts to equal retryMax")
		})
	}
}

func TestRetry_DefaultRetryMax(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 1000}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
	)

	require.Error(t, err, "Expected the operation to fail after the default number of attempts")
	assert.Equal(t, 3, op.calls, "Expected the default retryMax of 3 attempts")
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
		hqgoretrier.WithRetryMax(2),
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
	)

	require.ErrorIs(t, err, errAlternate, "Expected the error from the final attempt")
	assert.NotErrorIs(t, err, errTestOperation, "Expected the earlier error to be discarded")
}

func TestRetry_NotifierReceivesErrorAndBackoff(t *testing.T) {
	t.Parallel()

	type notification struct {
		err     error
		backoff time.Duration
	}

	var got []notification

	perAttempt := func(_, _ time.Duration, attempt int) time.Duration {
		return time.Duration(attempt) * time.Millisecond
	}

	op := &flakyOperation{failures: 2}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryBackoff(perAttempt),
		hqgoretrier.WithNotifier(func(err error, b time.Duration) {
			got = append(got, notification{err: err, backoff: b})
		}),
	)

	require.NoError(t, err, "Expected the operation to succeed after retries")
	require.Len(t, got, 2, "Expected the notifier to fire once per failed attempt")

	for i, n := range got {
		require.ErrorIs(t, n.err, errTestOperation, "Expected the notifier to receive the failure")
		assert.Equal(t, time.Duration(i+1)*time.Millisecond, n.backoff, "Expected the computed backoff")
	}
}

func TestRetry_BackoffReceivesConfiguredArgs(t *testing.T) {
	t.Parallel()

	type call struct {
		minDelay, maxDelay time.Duration
		attempt            int
	}

	var calls []call

	waitMin := 7 * time.Millisecond
	waitMax := 11 * time.Millisecond

	recording := func(minDelay, maxDelay time.Duration, attempt int) time.Duration {
		calls = append(calls, call{minDelay: minDelay, maxDelay: maxDelay, attempt: attempt})

		return 0
	}

	op := &flakyOperation{failures: 3}

	err := hqgoretrier.Retry(
		t.Context(),
		op.run,
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryWaitMin(waitMin),
		hqgoretrier.WithRetryWaitMax(waitMax),
		hqgoretrier.WithRetryBackoff(recording),
	)

	require.NoError(t, err, "Expected the operation to succeed after retries")
	require.Equal(t, []call{
		{minDelay: waitMin, maxDelay: waitMax, attempt: 1},
		{minDelay: waitMin, maxDelay: waitMax, attempt: 2},
		{minDelay: waitMin, maxDelay: waitMax, attempt: 3},
	}, calls, "Expected backoff to receive configured bounds and increasing attempt numbers")
}

func TestRetry_ZeroWaitMin(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 2}

	require.NotPanics(t, func() {
		err := hqgoretrier.Retry(
			t.Context(),
			op.run,
			hqgoretrier.WithRetryMax(5),
			hqgoretrier.WithRetryWaitMin(0),
			hqgoretrier.WithRetryWaitMax(0),
			hqgoretrier.WithRetryBackoff(backoff.Exponential()),
		)

		require.NoError(t, err, "Expected the operation to succeed after retries")
	})

	assert.Equal(t, 3, op.calls, "Expected the operation to be called 3 times")
}

func TestRetry_ContextCanceledBeforeStart(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	op := &flakyOperation{failures: 2}

	err := hqgoretrier.Retry(
		ctx,
		op.run,
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
	)

	require.ErrorIs(t, err, context.Canceled, "Expected a cancellation error")
	assert.Zero(t, op.calls, "Expected the operation never to run once the context is already canceled")
}

func TestRetry_ContextCanceledDuringWait(t *testing.T) {
	t.Parallel()

	cause := errors.New("canceled during wait")

	ctx, cancel := context.WithCancelCause(t.Context())
	defer cancel(nil)

	waiting := make(chan struct{})

	var once sync.Once

	longBackoff := func(_, _ time.Duration, _ int) time.Duration {
		once.Do(func() { close(waiting) })

		return time.Hour
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
		hqgoretrier.WithRetryMax(10),
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
		hqgoretrier.WithRetryMax(1000),
		hqgoretrier.WithRetryWaitMin(30*time.Millisecond),
		hqgoretrier.WithRetryWaitMax(100*time.Millisecond),
		hqgoretrier.WithRetryBackoff(backoff.Exponential()),
	)

	require.ErrorIs(t, err, context.DeadlineExceeded, "Expected a deadline-exceeded error")
	assert.LessOrEqual(t, op.calls, 2, "Expected few attempts before the deadline elapses")
}

func TestRetry_UnlimitedRetriesStopOnContext(t *testing.T) {
	t.Parallel()

	op := &flakyOperation{failures: 1_000_000}

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)

	defer cancel()

	err := hqgoretrier.Retry(
		ctx,
		op.run,
		hqgoretrier.WithRetryMax(0),
		hqgoretrier.WithRetryWaitMin(time.Millisecond),
		hqgoretrier.WithRetryWaitMax(2*time.Millisecond),
		hqgoretrier.WithRetryBackoff(backoff.Exponential()),
	)

	require.ErrorIs(t, err, context.DeadlineExceeded, "Expected the context deadline to stop unlimited retries")
	assert.Positive(t, op.calls, "Expected the operation to be attempted at least once")
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
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
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
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
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
		hqgoretrier.WithRetryMax(3),
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the last operation error")
	assert.Equal(t, 3, result, "Expected the data from the final attempt")
	assert.Equal(t, 3, calls, "Expected three attempts")
}

func TestRetryWithData_ContextCanceledReturnsZeroValue(t *testing.T) {
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
		hqgoretrier.WithRetryMax(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff()),
	)

	require.ErrorIs(t, err, context.Canceled, "Expected a cancellation error")
	assert.Zero(t, result, "Expected the zero value when canceled before any attempt")
	assert.Zero(t, calls, "Expected the operation never to run")
}
