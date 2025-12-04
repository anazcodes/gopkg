package retry

import (
	"context"
	"errors"
	"time"
)

// Config defines retry behavior.
type Config struct {
	MaxRetries int
	Delay      time.Duration
	Backoff    bool // if true, doubles delay each retry
}

// Do executes the given function with retry logic.
func Do(ctx context.Context, cfg Config, fn func() error) error {
	if cfg.MaxRetries <= 0 {
		cfg.MaxRetries = 1
	}

	delay := cfg.Delay
	var lastErr error

	for i := 0; i < cfg.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			lastErr = fn()
			if lastErr == nil {
				return nil
			}

			// Wait before retrying
			if i < cfg.MaxRetries-1 {
				time.Sleep(delay)
				if cfg.Backoff {
					delay *= 2
				}
			}
		}
	}

	return errors.New("max retries reached: " + lastErr.Error())
}
