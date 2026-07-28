package retrier

// permanentError wraps an error to mark it as non-retryable. Retry and RetryWithData detect it
// with errors.As, so the marker keeps working even when the operation wraps it further with
// fmt.Errorf and %w.
type permanentError struct {
	err error
}

// Error returns the wrapped error's message, keeping the marker invisible in output.
//
// Returns:
//   - (string): The message of the wrapped error.
func (e *permanentError) Error() (message string) {
	return e.err.Error()
}

// Unwrap exposes the wrapped error, so errors.Is and errors.As keep working on it.
//
// Returns:
//   - (error): The wrapped error.
func (e *permanentError) Unwrap() (err error) {
	return e.err
}

// Permanent wraps err to mark it as permanent: when an operation returns an error containing
// this marker, Retry and RetryWithData stop immediately instead of scheduling another attempt,
// and return the error to the caller.
//
// The marker is invisible to the caller: the returned error reads and matches exactly as the
// operation wrote it. The marker also survives further wrapping, so an operation may add context
// with fmt.Errorf and %w on top of a permanent error and it still stops the retry loop. Permanent
// errors are never retried, regardless of the WithRetryIf predicate.
//
// Typical candidates are failures that repeat identically on every attempt — a malformed request
// the API keeps rejecting, an unknown host, a missing record — anything another wait cannot fix.
//
// Parameters:
//   - err (error): The error to mark as permanent.
//
// Returns:
//   - (error): The marked error, or nil if err is nil.
func Permanent(err error) (perr error) {
	if err == nil {
		return nil
	}

	return &permanentError{err: err}
}
