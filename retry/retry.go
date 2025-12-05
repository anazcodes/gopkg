package retry

import (
	"context"
	"errors"
	"time"
)

// Client defines retry behavior.
type Client struct {
	MaxRetries int
	Delay      time.Duration
	Backoff    bool // if true, doubles delay each retry
}

// Do executes the given function with retry logic.
func (c *Client) Do(ctx context.Context, fn func() error) error {
	if c.MaxRetries <= 0 {
		c.MaxRetries = 1
	}

	delay := c.Delay
	var lastErr error

	for i := 0; i < c.MaxRetries; i++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			lastErr = fn()
			if lastErr == nil {
				return nil
			}

			// Wait before retrying
			if i < c.MaxRetries-1 {
				time.Sleep(delay)
				if c.Backoff {
					delay *= 2
				}
			}
		}
	}

	return errors.New("max retries reached: " + lastErr.Error())
}
