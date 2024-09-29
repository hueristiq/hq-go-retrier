package hqgoretry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hueristiq/hqgoretry"
	"github.com/hueristiq/hqgoretry/backoff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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
	ctx := context.Background()

	err := hqgoretry.Retry(ctx, mockOp.Operation,
		hqgoretry.WithMaxRetries(5),
		hqgoretry.WithMinDelay(10*time.Millisecond),
		hqgoretry.WithMaxDelay(50*time.Millisecond),
		hqgoretry.WithBackoff(backoff.Exponential()))

	require.NoError(t, err, "Expected operation to succeed after retries")
	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetry_MaxRetriesExceeded(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 10} // Will fail more times than the allowed retries
	ctx := context.Background()

	err := hqgoretry.Retry(ctx, mockOp.Operation,
		hqgoretry.WithMaxRetries(3),
		hqgoretry.WithMinDelay(10*time.Millisecond),
		hqgoretry.WithMaxDelay(50*time.Millisecond),
		hqgoretry.WithBackoff(backoff.Exponential()))

	require.Error(t, err, "Expected operation to fail after retries")
	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetryWithContext_Timeout(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 10}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	err := hqgoretry.Retry(ctx, mockOp.Operation,
		hqgoretry.WithMaxRetries(5),
		hqgoretry.WithMinDelay(30*time.Millisecond),
		hqgoretry.WithMaxDelay(100*time.Millisecond),
		hqgoretry.WithBackoff(backoff.Exponential()))

	require.Error(t, err, "Expected operation to fail due to context timeout")
	require.ErrorIs(t, err, context.DeadlineExceeded, "Expected timeout error")
	assert.LessOrEqual(t, mockOp.callCount, 2, "Expected the operation to be called less than the max retries due to timeout")
}

func TestRetryWithData_Success(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}
	ctx := context.Background()

	operationWithData := func() (int, error) {
		if mockOp.callCount < 2 {
			mockOp.callCount++

			return 0, errTestOperation
		}

		return 42, nil
	}

	result, err := hqgoretry.RetryWithData(ctx, operationWithData,
		hqgoretry.WithMaxRetries(5),
		hqgoretry.WithMinDelay(10*time.Millisecond),
		hqgoretry.WithMaxDelay(50*time.Millisecond),
		hqgoretry.WithBackoff(backoff.Exponential()))

	require.NoError(t, err, "Expected operation to succeed after retries")
	assert.Equal(t, 42, result, "Expected operation result to be 42")
}

func TestRetryWithDecorrelatedJitter(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}
	ctx := context.Background()

	err := hqgoretry.Retry(ctx, mockOp.Operation,
		hqgoretry.WithMaxRetries(5),
		hqgoretry.WithMinDelay(10*time.Millisecond),
		hqgoretry.WithMaxDelay(50*time.Millisecond),
		hqgoretry.WithBackoff(backoff.ExponentialWithDecorrelatedJitter()))

	require.NoError(t, err, "Expected operation to succeed after retries with decorrelated jitter")
	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetry_FullJitter(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}
	ctx := context.Background()

	err := hqgoretry.Retry(ctx, mockOp.Operation,
		hqgoretry.WithMaxRetries(5),
		hqgoretry.WithMinDelay(10*time.Millisecond),
		hqgoretry.WithMaxDelay(50*time.Millisecond),
		hqgoretry.WithBackoff(backoff.ExponentialWithFullJitter()))

	require.NoError(t, err, "Expected operation to succeed after retries with full jitter")
	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetry_EqualJitter(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}
	ctx := context.Background()

	err := hqgoretry.Retry(ctx, mockOp.Operation,
		hqgoretry.WithMaxRetries(5),
		hqgoretry.WithMinDelay(10*time.Millisecond),
		hqgoretry.WithMaxDelay(50*time.Millisecond),
		hqgoretry.WithBackoff(backoff.ExponentialWithEqualJitter()))

	require.NoError(t, err, "Expected operation to succeed after retries with equal jitter")
	assert.Equal(t, 3, mockOp.callCount, "Expected the operation to be called 3 times")
}

func TestRetry_ContextCanceled(t *testing.T) {
	t.Parallel()

	mockOp := &mockOperation{failureCount: 2}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := hqgoretry.Retry(ctx, mockOp.Operation,
		hqgoretry.WithMaxRetries(5),
		hqgoretry.WithMinDelay(10*time.Millisecond),
		hqgoretry.WithMaxDelay(50*time.Millisecond),
		hqgoretry.WithBackoff(backoff.Exponential()))

	require.Error(t, err, "Expected operation to fail due to canceled context")
	require.ErrorIs(t, err, context.Canceled, "Expected timeout error")
}
