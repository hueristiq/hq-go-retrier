package hqgoretry

import (
	"context"
	"time"

	"github.com/hueristiq/hqgoretry/backoff"
)

type Operation func() (err error)

func (o Operation) withEmptyData() OperationWithData[struct{}] {
	return func() (struct{}, error) {
		return struct{}{}, o()
	}
}

type OperationWithData[T any] func() (data T, err error)

type Configuration struct {
	maxRetries int
	minDelay   time.Duration
	maxDelay   time.Duration
	backoff    backoff.Backoff
}

type Option func(*Configuration)

func WithMaxRetries(retries int) Option {
	return func(c *Configuration) {
		c.maxRetries = retries
	}
}

func WithMaxDelay(delay time.Duration) Option {
	return func(c *Configuration) {
		c.maxDelay = delay
	}
}

func WithMinDelay(delay time.Duration) Option {
	return func(c *Configuration) {
		c.minDelay = delay
	}
}

func WithBackoff(strategy backoff.Backoff) Option {
	return func(c *Configuration) {
		c.backoff = strategy
	}
}

func Retry(ctx context.Context, operation Operation, opts ...Option) (err error) {
	_, err = RetryWithData(ctx, operation.withEmptyData(), opts...)

	return
}

func RetryWithData[T any](ctx context.Context, operation OperationWithData[T], opts ...Option) (result T, err error) {
	cfg := &Configuration{
		maxRetries: 3,
		maxDelay:   1000 * time.Millisecond,
		minDelay:   100 * time.Millisecond,
		backoff:    backoff.Exponential(),
	}

	for _, opt := range opts {
		opt(cfg)
	}

	for attempt := range cfg.maxRetries {
		select {
		case <-ctx.Done():
			err = ctx.Err()

			return
		default:
			result, err = operation()
			if err == nil {
				return
			}

			b := cfg.backoff(cfg.minDelay, cfg.maxDelay, attempt)

			ticker := time.NewTicker(b)

			select {
			case <-ticker.C:
				ticker.Stop()
			case <-ctx.Done():
				ticker.Stop()

				err = ctx.Err()

				return
			}
		}
	}

	return
}
