package retrier

import (
	"context"
	"time"

	"github.com/hueristiq/hq-go-retrier/backoff"
)

type Operation func() (err error)

func (o Operation) withEmptyData() (operationWithData OperationWithData[struct{}]) {
	operationWithData = func() (struct{}, error) {
		return struct{}{}, o()
	}

	return
}

type OperationWithData[T any] func() (data T, err error)

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
