package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"
	"sync"
)

type ctxKey string

const (
	slogFields ctxKey = "slog_fields"
)

const callerDepth = 9

type LoggerInterface interface {
	Debug(msg string, args ...any)
	DebugContext(ctx context.Context, msg string, args ...any)

	Info(msg string, args ...any)
	InfoContext(ctx context.Context, msg string, args ...any)

	Warn(msg string, args ...any)
	WarnContext(ctx context.Context, msg string, args ...any)

	Error(msg string, args ...any)
	ErrorContext(ctx context.Context, msg string, args ...any)

	With(args ...any) LoggerInterface
	WithGroup(name string) LoggerInterface

	Log(ctx context.Context, level slog.Level, msg string, args ...any)
}

type Logger struct {
	log *slog.Logger
}

type contextHandler struct {
	handler slog.Handler
}

var _ LoggerInterface = (*Logger)(nil)

// InitDefault инициализирует глобальный логгер
func InitDefault(cfg *Config) error {
	logHandler, err := newHandler(cfg)
	if err != nil {
		return err
	}

	l := slog.New(logHandler)
	slog.SetDefault(l)
	return nil
}

// Default получение логгера по умолчанию
func Default() LoggerInterface {
	return &Logger{log: slog.Default()}
}

// New получение нового логгера
func New(cfg *Config) (LoggerInterface, error) {
	logHandler, err := newHandler(cfg)
	if err != nil {
		return nil, err
	}
	return &Logger{log: slog.New(logHandler)}, nil
}

func (l *Logger) Debug(msg string, args ...any) {
	l.log.Debug(msg, args...)
}

func (l *Logger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.log.DebugContext(ctx, msg, args...)
}

func (l *Logger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l *Logger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.log.InfoContext(ctx, msg, args...)
}

func (l *Logger) Warn(msg string, args ...any) {
	l.log.Warn(msg, args...)
}

func (l *Logger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.log.WarnContext(ctx, msg, args...)
}

func (l *Logger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.log.ErrorContext(ctx, msg, args...)
}

func (l *Logger) With(args ...any) LoggerInterface {
	return &Logger{log: l.log.With(args...)}
}

func (l *Logger) WithGroup(name string) LoggerInterface {
	return &Logger{log: l.log.WithGroup(name)}
}

func (l *Logger) Log(ctx context.Context, level slog.Level, msg string, args ...any) {
	l.log.Log(ctx, level, msg, args...)
}

func newContextHandler(handler slog.Handler) *contextHandler {
	return &contextHandler{handler: handler}
}

func (ch *contextHandler) Enabled(ctx context.Context, rec slog.Level) bool {
	return ch.handler.Enabled(ctx, rec)
}

func (ch *contextHandler) Handle(ctx context.Context, r slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).(*sync.Map); ok {
		attrs.Range(func(key, val any) bool {
			if keyString, ok := key.(string); ok {
				r.AddAttrs(slog.Any(keyString, val))
			}
			return true
		})
	}
	return ch.handler.Handle(ctx, r)
}

func (ch *contextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &contextHandler{handler: ch.handler.WithAttrs(attrs)}
}

func (ch *contextHandler) WithGroup(name string) slog.Handler {
	return &contextHandler{handler: ch.handler.WithGroup(name)}
}

// AppendCtx добавление пары ключ/значение в контекст
func AppendCtx(parent context.Context, key string, val any) context.Context {
	if parent == nil {
		parent = context.Background()
	}

	if v, ok := parent.Value(slogFields).(*sync.Map); ok {
		v.Store(key, val)
		return context.WithValue(parent, slogFields, v)
	}
	v := &sync.Map{}
	v.Store(key, val)
	return context.WithValue(parent, slogFields, v)
}

func replaceAttrFunc(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.SourceKey {
		source := a.Value.Any().(*slog.Source)
		pc := make([]uintptr, 1)
		runtime.Callers(callerDepth, pc)
		fs := runtime.CallersFrames([]uintptr{pc[0]})
		f, _ := fs.Next()

		source.File = f.File
		source.Function = f.Function
		source.Line = f.Line
	}
	return a
}

func newHandler(cfg *Config) (slog.Handler, error) {
	var logWriter io.Writer
	var logLevel slog.Level

	if cfg.Path == "" {
		logWriter = os.Stdout
	} else {
		logFile, err := os.OpenFile(cfg.Path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
		if err != nil {
			return nil, fmt.Errorf("failed to open log file %s: %w", cfg.Path, err)
		}
		logWriter = logFile
	}

	switch cfg.Level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelError
	}

	logOptions := &slog.HandlerOptions{
		Level:       logLevel,
		AddSource:   true,
		ReplaceAttr: replaceAttrFunc,
	}
	return newContextHandler(slog.NewJSONHandler(logWriter, logOptions)), nil
}
