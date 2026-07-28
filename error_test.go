package retrier_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	hqgoretrier "github.com/hueristiq/hq-lib-retrier-go"
)

func TestPermanent_ErrorMessageIsUnchanged(t *testing.T) {
	t.Parallel()

	err := hqgoretrier.Permanent(errTestOperation)

	assert.Equal(t, errTestOperation.Error(), err.Error(), "Expected the marker to be invisible in the error message")
}

func TestPermanent_UnwrapExposesOriginal(t *testing.T) {
	t.Parallel()

	err := hqgoretrier.Permanent(errTestOperation)

	require.ErrorIs(t, err, errTestOperation, "Expected errors.Is to match the wrapped error")
	assert.Equal(t, errTestOperation, errors.Unwrap(err), "Expected Unwrap to return the original error")
}

func TestPermanent_DoubleMarked(t *testing.T) {
	t.Parallel()

	calls := 0

	op := func() error {
		calls++

		return hqgoretrier.Permanent(hqgoretrier.Permanent(errTestOperation))
	}

	err := hqgoretrier.Retry(
		t.Context(),
		op,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the permanent error")
	assert.Equal(t, 1, calls, "Expected no retries even when the marker is applied twice")
}

func TestRetry_PermanentErrorStopsRetries(t *testing.T) {
	t.Parallel()

	t.Run("returned directly", func(t *testing.T) {
		t.Parallel()

		calls := 0

		op := func() error {
			calls++

			return hqgoretrier.Permanent(errTestOperation)
		}

		err := hqgoretrier.Retry(
			t.Context(),
			op,
			hqgoretrier.WithMaxAttempts(5),
			hqgoretrier.WithRetryBackoff(noWaitBackoff),
		)

		require.ErrorIs(t, err, errTestOperation, "Expected the permanent error")
		assert.Equal(t, 1, calls, "Expected no retries after a permanent error")
	})

	t.Run("wrapped in operation context", func(t *testing.T) {
		t.Parallel()

		calls := 0

		op := func() error {
			calls++

			return fmt.Errorf("fetching data: %w", hqgoretrier.Permanent(errTestOperation))
		}

		err := hqgoretrier.Retry(
			t.Context(),
			op,
			hqgoretrier.WithMaxAttempts(5),
			hqgoretrier.WithRetryBackoff(noWaitBackoff),
		)

		require.ErrorIs(t, err, errTestOperation, "Expected errors.Is to reach the permanent error")
		require.ErrorContains(t, err, "fetching data", "Expected the operation's wrap context to be preserved")
		assert.Equal(t, 1, calls, "Expected no retries after a permanent error")
	})

	t.Run("wrapped twice", func(t *testing.T) {
		t.Parallel()

		calls := 0

		op := func() error {
			calls++

			return fmt.Errorf("outer: %w", fmt.Errorf("inner: %w", hqgoretrier.Permanent(errTestOperation)))
		}

		err := hqgoretrier.Retry(
			t.Context(),
			op,
			hqgoretrier.WithMaxAttempts(5),
			hqgoretrier.WithRetryBackoff(noWaitBackoff),
		)

		require.ErrorIs(t, err, errTestOperation, "Expected errors.Is to reach the permanent error through two wraps")
		assert.Equal(t, 1, calls, "Expected no retries after a permanent error")
	})

	t.Run("nil error", func(t *testing.T) {
		t.Parallel()

		assert.NoError(t, hqgoretrier.Permanent(nil), "Expected Permanent(nil) to return nil")
	})
}

func TestRetry_PermanentErrorWinsOverRetryIf(t *testing.T) {
	t.Parallel()

	calls := 0

	op := func() error {
		calls++

		return hqgoretrier.Permanent(errTestOperation)
	}

	err := hqgoretrier.Retry(
		t.Context(),
		op,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
		hqgoretrier.WithRetryIf(func(error) bool { return true }),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the permanent error")
	assert.Equal(t, 1, calls, "Expected permanent errors never to be retried")
}

func TestRetryWithData_PermanentErrorReturnsLastResult(t *testing.T) {
	t.Parallel()

	calls := 0

	op := func() (int, error) {
		calls++

		return calls, hqgoretrier.Permanent(errTestOperation)
	}

	result, err := hqgoretrier.RetryWithData(
		t.Context(),
		op,
		hqgoretrier.WithMaxAttempts(5),
		hqgoretrier.WithRetryBackoff(noWaitBackoff),
	)

	require.ErrorIs(t, err, errTestOperation, "Expected the permanent error")
	assert.Equal(t, 1, result, "Expected the data from the only attempt")
	assert.Equal(t, 1, calls, "Expected no retries after a permanent error")
}
