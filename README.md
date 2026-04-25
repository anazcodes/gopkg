# gopkg

[![Go Reference](https://pkg.go.dev/badge/github.com/anazcodes/gopkg.svg)](https://pkg.go.dev/github.com/anazcodes/gopkg)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

`gopkg` is a collection of reusable, robust Go packages and utilities designed to accelerate backend development.

## Packages

- `goslack`: Utilities for interacting with Slack.
- `httpclient`: A robust HTTP client wrapper.
- `logger`: Structured logging configuration and hooks (e.g., using `zerolog`).
- `otel`: OpenTelemetry integrations for tracing and metrics.
- `pgxtype`: Extensions and utilities for working with PostgreSQL via `pgx`.
- `retry`: Robust retry logic with exponential backoff support.

## Installation

```bash
go get github.com/anazcodes/gopkg
```

## Example: Retry

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/anazcodes/gopkg/retry"
)

func main() {
	r := retry.Retryer{
		MaxRetries: 3,
		Delay:      1 * time.Second,
		Backoff:    true,
	}

	err := r.Do(context.Background(), func() error {
		// Your logic here
		return fmt.Errorf("temporary error")
	})

	if err != nil {
		fmt.Println("Failed after retries:", err)
	}
}
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
