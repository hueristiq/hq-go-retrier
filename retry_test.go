package retrier_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.source.hueristiq.com/retrier"
	"go.source.hueristiq.com/retrier/backoff"
)

var errTestOperation = errors.New("operation failed")

// Mock operation that will fail a given number of times before succeeding.
type mockOperation struct {
	mock.Mock

	failureCount int
	callCount    int
}

func (m *mockOperation) Operation() error {
	m.callCount++

	if m.callCount <= m.failureCount {
		return errTestOperation
	}

	return nil
}

func TestRetry_SuccessAfterFailures(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2} // Fail twice, then succeed
	ctx := t.Context()

	err := retrier.Retry(ctx, mockOp.Operation,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(10*time.Millisecond),
		retrier.WithRetryWaitMax(50*time.Millisecond),
		retrier.WithRetryBackoff(backoff.Exponential()))

	require.NoError(t, err, "Expected operation to succeed after retries")

	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetry_MaxRetriesExceeded(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 10} // Will fail more times than the allowed retries
	ctx := t.Context()

	err := retrier.Retry(ctx, mockOp.Operation,
		retrier.WithRetryMax(3),
		retrier.WithRetryWaitMin(10*time.Millisecond),
		retrier.WithRetryWaitMax(50*time.Millisecond),
		retrier.WithRetryBackoff(backoff.Exponential()))

	require.Error(t, err, "Expected operation to fail after retries")

	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetryWithContext_Timeout(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 10}

	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)

	defer cancel()

	err := retrier.Retry(ctx, mockOp.Operation,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(30*time.Millisecond),
		retrier.WithRetryWaitMax(100*time.Millisecond),
		retrier.WithRetryBackoff(backoff.Exponential()))

	require.Error(t, err, "Expected operation to fail due to context timeout")
	require.ErrorIs(t, err, context.DeadlineExceeded, "Expected timeout error")

	assert.LessOrEqual(t, mockOp.callCount, 2, "Expected the operation to be called less than the max retries due to timeout")
}

func TestRetryWithData_Success(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}
	ctx := t.Context()

	operationWithData := func() (int, error) {
		if mockOp.callCount < 2 {
			mockOp.callCount++

			return 0, errTestOperation
		}

		return 42, nil
	}

	result, err := retrier.RetryWithData(ctx, operationWithData,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(10*time.Millisecond),
		retrier.WithRetryWaitMax(50*time.Millisecond),
		retrier.WithRetryBackoff(backoff.Exponential()))

	require.NoError(t, err, "Expected operation to succeed after retries")

	assert.Equal(t, 42, result, "Expected operation result to be 42")
}

func TestRetryWithDecorrelatedJitter(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}
	ctx := t.Context()

	err := retrier.Retry(ctx, mockOp.Operation,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(10*time.Millisecond),
		retrier.WithRetryWaitMax(50*time.Millisecond),
		retrier.WithRetryBackoff(backoff.ExponentialWithDecorrelatedJitter()))

	require.NoError(t, err, "Expected operation to succeed after retries with decorrelated jitter")

	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetry_FullJitter(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}
	ctx := t.Context()

	err := retrier.Retry(ctx, mockOp.Operation,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(10*time.Millisecond),
		retrier.WithRetryWaitMax(50*time.Millisecond),
		retrier.WithRetryBackoff(backoff.ExponentialWithFullJitter()))

	require.NoError(t, err, "Expected operation to succeed after retries with full jitter")

	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetry_EqualJitter(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}
	ctx := t.Context()

	err := retrier.Retry(ctx, mockOp.Operation,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(10*time.Millisecond),
		retrier.WithRetryWaitMax(50*time.Millisecond),
		retrier.WithRetryBackoff(backoff.ExponentialWithEqualJitter()))

	require.NoError(t, err, "Expected operation to succeed after retries with equal jitter")

	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetry_ContextCanceled(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}

	ctx, cancel := context.WithCancel(t.Context())

	cancel()

	err := retrier.Retry(ctx, mockOp.Operation,
		retrier.WithRetryMax(5),
		retrier.WithRetryWaitMin(10*time.Millisecond),
		retrier.WithRetryWaitMax(50*time.Millisecond),
		retrier.WithRetryBackoff(backoff.Exponential()))

	require.Error(t, err, "Expected operation to fail due to canceled context")
	require.ErrorIs(t, err, context.Canceled, "Expected timeout error")
}
