package log

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"
)

var LevelFatal slog.Level = slog.Level(12)
var LevelTrace slog.Level = slog.Level(-8)
var LevelUnknown slog.Level = slog.Level(-99)

// // PrettyHandler makes the log lines look nice instead of pure key-value pairs
// type PrettyHandler struct {
// 	handler slog.Handler
// }

// func NewPrettyHandler(handler slog.Handler) *PrettyHandler {
// 	if prettyHandler, ok := handler.(*PrettyHandler); ok {
// 		handler = prettyHandler.Handler()
// 	}
// 	return &PrettyHandler{handler}
// }

// // Enabled implements Handler.Enabled by reporting whether
// // level is at least as large as h's level.
// func (h *PrettyHandler) Enabled(ctx context.Context, level slog.Level) bool {
// 	return h.handler.Enabled(ctx, level)
// }

// // Handle implements Handler.Handle.
// func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
// 	return h.handler.Handle(ctx, r)
//     // logLine := fmt.Sprintf()
//     // fmt.Printf()
// }

// // WithAttrs implements Handler.WithAttrs.
// func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
// 	return h.handler.WithAttrs(attrs)
// }

// // WithGroup implements Handler.WithGroup.
// func (h *PrettyHandler) WithGroup(name string) slog.Handler {
// 	return h.handler.WithGroup(name)
// }

// // Handler returns the Handler wrapped by h.
// func (h *PrettyHandler) Handler() slog.Handler {
// 	return h.handler
// }

// newLogger creates the logger with the desired level
func newLogger(level slog.Level) *slog.Logger {
	// custom handler so FATAL and TRACE level strings are represented in the log
	textHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: level,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				switch level {
				case LevelFatal:
					a.Value = slog.StringValue("FATAL")
				case LevelTrace:
					a.Value = slog.StringValue("TRACE")
				}
			}
			return a
		},
	})

	// customHandler := NewPrettyHandler(textHandler)
	return slog.New(textHandler)
}

// ParseLevel parses the log level string to an slog.Level
func ParseLevel(str string) (slog.Level, error) {
	switch strings.ToUpper(str) {
	case "ERROR":
		return slog.LevelError, nil
	case "WARN":
		return slog.LevelWarn, nil
	case "INFO":
		return slog.LevelInfo, nil
	case "DEBUG":
		return slog.LevelDebug, nil
	case "TRACE":
		// slog has no trace level by default
		return LevelTrace, nil
	}

	return LevelUnknown, fmt.Errorf("cannot parse log level: %v", str)
}

// SetLevel sets the log level from a provided slog.Level
func SetLevel(level slog.Level) {
	logger := newLogger(level)
	slog.SetDefault(logger)
}

// Debug wraps the slog.Debug method
func Debug(msg string) {
	slog.LogAttrs(context.TODO(), slog.LevelDebug, msg)
}

// Debugf wraps the slog.Debug method
// and passes the arguments to fmt.Sprintf
func Debugf(msg string, args ...any) {
	slog.LogAttrs(context.TODO(), slog.LevelDebug, fmt.Sprintf(msg, args...))
}

// Error wraps the slog.Error method
func Error(msg string, args ...any) {
	slog.LogAttrs(context.TODO(), slog.LevelError, msg)
}

// Errorf wraps the slog.Error method
// and passes the arguments to fmt.Sprintf
func Errorf(msg string, args ...any) {
	slog.LogAttrs(context.TODO(), slog.LevelError, fmt.Sprintf(msg, args...))
}

// Fatal wraps the slog.Log method with a custom log level,
// and mimics go's log.Fatal which calls os.Exit(1)
// https://pkg.go.dev/log#Fatal
func Fatal(msg string) {
	slog.LogAttrs(context.TODO(), LevelFatal, msg)
	os.Exit(1)
}

// Fatal wraps the slog.Log method with a custom log level,
// passes the arguments to fmt.Sprintf,
// and mimics go's log.Fatal which calls os.Exit(1)
// https://pkg.go.dev/log#Fatal
func Fatalf(msg string, args ...any) {
	slog.LogAttrs(context.TODO(), LevelFatal, fmt.Sprintf(msg, args...))
	os.Exit(1)
}

// Info wraps the slog.Info method
func Info(msg string) {
	slog.LogAttrs(context.TODO(), slog.LevelInfo, msg)
}

// Infof wraps the slog.Info method
// and passes the arguments to fmt.Sprintf
func Infof(msg string, args ...any) {
	slog.LogAttrs(context.TODO(), slog.LevelInfo, fmt.Sprintf(msg, args...))
}

// Trace wraps the slog.Log method with a custom log level,
func Trace(msg string) {
	slog.LogAttrs(context.TODO(), LevelTrace, msg)
}

// Tracef wraps the slog.Log method with a custom log level
// and passes the arguments to fmt.Sprintf
func Tracef(msg string, args ...any) {
	slog.LogAttrs(context.TODO(), LevelTrace, fmt.Sprintf(msg, args...))

}

// Warn wraps the slog.Warn method
func Warn(msg string) {
	slog.LogAttrs(context.TODO(), slog.LevelWarn, msg)
}

// Warnf wraps the slog.Warn method
func Warnf(msg string, args ...any) {
	slog.LogAttrs(context.TODO(), slog.LevelWarn, fmt.Sprintf(msg, args...))

}
