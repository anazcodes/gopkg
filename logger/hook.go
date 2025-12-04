package logger

import (
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/Turinix/gopkg/otel"
	"github.com/rs/zerolog"
)

type TracingHook struct {
	// MaxDepth is the depth of the caller file-path need to log. Depth calculates in reverse from the caller file.
	//
	// Default: 0
	MaxDepth int
}

func (h TracingHook) Run(e *zerolog.Event, _ zerolog.Level, _ string) {
	ctx := e.GetCtx()
	traceID := otel.GetTraceID(ctx)

	e.Str("traceId", traceID).Str("file", callerWithDepth(h.MaxDepth))
}

// callerWithDepth returns caller file path with line number.
// The depth is calculated in reverse from the caller file.
func callerWithDepth(maxDepth int) string {
	const skipFrames = 5

	_, file, line, ok := runtime.Caller(skipFrames)
	if !ok {
		return "unknown:0"
	}

	cleanPath := filepath.Clean(file)
	components := strings.Split(cleanPath, string(filepath.Separator))

	start := len(components) - maxDepth
	if start < 0 {
		start = 0
	}
	selected := components[start:]

	callerPath := strings.Join(selected, "/")

	return callerPath + ":" + strconv.Itoa(line)
}
