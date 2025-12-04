package logger

type OptFunc func(*Opts)

func defaultOptions() *Opts {
	return &Opts{
		Filename:    "log/log.json",
		MaxSize:     512,
		MaxAge:      3,
		MaxBackups:  10,
		LocalTime:   false,
		Compress:    false,
		CallerDepth: 3,
		StdOut:      false,
		Level:       -1,
	}
}

type Opts struct {
	Filename    string
	MaxSize     int
	MaxAge      int
	MaxBackups  int
	LocalTime   bool
	Compress    bool
	CallerDepth int
	StdOut      bool
	Level       int8
}

func WithFileName(f string) OptFunc {
	return func(o *Opts) {
		o.Filename = f
	}
}

func WithMaxSize(size int) OptFunc {
	return func(o *Opts) {
		o.MaxSize = size
	}
}

func WithMaxBackups(n int) OptFunc {
	return func(o *Opts) {
		o.MaxBackups = n
	}
}

func WithMaxAge(days int) OptFunc {
	return func(o *Opts) {
		o.MaxSize = days
	}
}

func WithLocalTime() OptFunc {
	return func(o *Opts) {
		o.LocalTime = true
	}
}

func WithCompress() OptFunc {
	return func(o *Opts) {
		o.Compress = true
	}
}

func WithCallerDepth(n int) OptFunc {
	return func(o *Opts) {
		o.CallerDepth = n
	}
}

func WithStdOut() OptFunc {
	return func(o *Opts) {
		o.StdOut = true
	}
}

func WithLevel(lvl int8) OptFunc {
	return func(o *Opts) {
		o.Level = lvl
	}
}
