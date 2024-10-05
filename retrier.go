package retrier

import (
	"time"

	"github.com/hueristiq/hq-go-retrier/backoff"
)

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
