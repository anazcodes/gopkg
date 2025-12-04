package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type logger struct {
	logger zerolog.Logger
}

// NewLogger returns a new logger.
// If opts is empty, the default options will be used.
func NewLogger(opts ...OptFunc) Logger {
	opt := defaultOptions()

	for _, fn := range opts {
		fn(opt)
	}

	return &logger{
		newLogger(opt),
	}
}

// newLogger configure and returns the newLogger.
func newLogger(opt *Opts) zerolog.Logger {
	roller := &lumberjack.Logger{
		Filename:   opt.Filename,
		MaxSize:    opt.MaxSize,
		MaxBackups: opt.MaxBackups,
		MaxAge:     opt.MaxAge,
		Compress:   opt.Compress,
		LocalTime:  opt.LocalTime,
	}

	ioWriters := []io.Writer{roller}

	if opt.StdOut {
		stdOut := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: "2006-01-02 03:04:05 PM",
			NoColor:    false,
		}

		ioWriters = append(ioWriters, stdOut)
	}

	zerolog.TimeFieldFormat = "2006-01-02 03:04:05 PM"
	zerolog.TimestampFunc = func() time.Time {
		t := time.Now()

		location := "Asia/Kolkata"
		loc, err := time.LoadLocation(location)
		if err != nil {
			log.Fatalf("failed logger.time.LoadLocation: %s", location)
		}

		return t.In(loc)
	}

	zerolog.SetGlobalLevel(zerolog.Level(opt.Level))

	writer := zerolog.MultiLevelWriter(ioWriters...)
	logger := zerolog.New(writer).With().Timestamp().Logger()

	return logger.Hook(TracingHook{MaxDepth: opt.CallerDepth})
}

type Logger interface {
	// Println implements the functionality of fmt.Println with level debug.
	Println(v ...any)
	// Trace logs in level trace with a message and expects a key value pair in args.
	Trace(ctx context.Context, args ...any)
	// Tracef logs in level Tracef with a message and supports format specifiers.
	Tracef(ctx context.Context, format string, v ...any)
	// Debug logs in level debug with a message and expects a key value pair in args.
	Debug(ctx context.Context, args ...any)
	// Debugf logs in level debug with a message and supports format specifiers.
	Debugf(ctx context.Context, format string, v ...any)
	// Error logs in level debug with a message and expects a key value pair in args.
	Error(ctx context.Context, args ...any)
	// Errorf logs in level debug with a message and supports format specifiers.
	Errorf(ctx context.Context, format string, v ...any)
	// ReturnErr logs in level error and returns the wrapped error.
	ReturnErr(ctx context.Context, msg string, err error) error
	// Fatal logs in level debug with a message and expects a key value pair in args.
	Fatal(ctx context.Context, args ...any)
	// Fatalf logs in level debug with a message and supports format specifiers.
	Fatalf(ctx context.Context, format string, v ...any)
	// Info logs in level debug with a message and expects a key value pair in args.
	Info(ctx context.Context, args ...any)
	// Infof logs in level debug with a message and supports format specifiers.
	Infof(ctx context.Context, format string, v ...any)
	// Warn logs in level debug with a message and expects a key value pair in args.
	Warn(ctx context.Context, args ...any)
	// Warnf logs in level debug with a message and supports format specifiers.
	Warnf(ctx context.Context, format string, v ...any)
}

func (l *logger) Println(v ...any) {
	l.logger.Println(v...)
}
func (l *logger) Trace(ctx context.Context, args ...any) {
	l.logger.Debug().Ctx(ctx).Fields(args).Send()
}

func (l *logger) Tracef(ctx context.Context, format string, v ...any) {
	l.logger.Debug().Ctx(ctx).Msgf(format, v...)
}

func (l *logger) Debug(ctx context.Context, args ...any) {
	l.logger.Debug().Ctx(ctx).Fields(args).Send()
}

func (l *logger) Debugf(ctx context.Context, format string, v ...any) {
	l.logger.Debug().Ctx(ctx).Msgf(format, v...)
}
func (l *logger) Info(ctx context.Context, args ...any) {
	l.logger.Info().Ctx(ctx).Fields(args).Send()
}

func (l *logger) Infof(ctx context.Context, format string, v ...any) {
	l.logger.Info().Ctx(ctx).Msgf(format, v...)
}

func (l *logger) Warn(ctx context.Context, args ...any) {
	l.logger.Warn().Ctx(ctx).Fields(args).Send()
}
func (l *logger) Warnf(ctx context.Context, format string, v ...any) {
	l.logger.Warn().Ctx(ctx).Msgf(format, v...)
}

func (l *logger) Error(ctx context.Context, args ...any) {
	l.logger.Error().Ctx(ctx).Fields(args).Send()
}

func (l *logger) Errorf(ctx context.Context, format string, v ...any) {
	l.logger.Error().Ctx(ctx).Msgf(format, v...)
}

func (l *logger) ReturnErr(ctx context.Context, msg string, err error) error {
	l.logger.Error().Ctx(ctx).Err(err).Msg(msg)

	return fmt.Errorf(msg+": %w", err)
}

func (l *logger) Fatal(ctx context.Context, args ...any) {
	l.logger.Fatal().Ctx(ctx).Fields(args).Send()
}

func (l *logger) Fatalf(ctx context.Context, format string, v ...any) {
	l.logger.Fatal().Ctx(ctx).Msgf(format, v...)
}
