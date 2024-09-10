package btclog

import (
	"context"
	"fmt"
	"log/slog"
)

// Handler wraps the slog.Handler interface with a few more methods that we
// need in order to satisfy the Logger interface.
type Handler interface {
	slog.Handler

	// Level returns the current logging level of the Handler.
	Level() Level

	// SetLevel changes the logging level of the Handler to the passed
	// level.
	SetLevel(level Level)
}

// sLogger is an implementation of Logger backed by a structured sLogger.
type sLogger struct {
	Handler
	logger *slog.Logger
}

// NewSLogger constructs a new structured logger from the given Handler.
func NewSLogger(handler Handler) Logger {
	return &sLogger{
		Handler: handler,
		logger:  slog.New(handler),
	}
}

// Tracef formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelTrace.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Tracef(format string, params ...any) {
	l.logger.Log(context.Background(), slog.Level(LevelTrace),
		fmt.Sprintf(format, params...))
}

// Debugf formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelDebug.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Debugf(format string, params ...any) {
	l.logger.Log(context.Background(), slog.Level(LevelDebug),
		fmt.Sprintf(format, params...))
}

// Infof formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelInfo.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Infof(format string, params ...any) {
	l.logger.Log(context.Background(), slog.Level(LevelInfo),
		fmt.Sprintf(format, params...))
}

// Warnf formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelWarn.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Warnf(format string, params ...any) {
	l.logger.Log(context.Background(), slog.Level(LevelWarn),
		fmt.Sprintf(format, params...))
}

// Errorf formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelError.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Errorf(format string, params ...any) {
	l.logger.Log(context.Background(), slog.Level(LevelError),
		fmt.Sprintf(format, params...))
}

// Criticalf formats message according to format specifier, prepends the prefix as
// necessary, and writes to log with LevelCritical.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Criticalf(format string, params ...any) {
	l.logger.Log(context.Background(), slog.Level(LevelCritical),
		fmt.Sprintf(format, params...))
}

// Trace formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelTrace.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Trace(v ...interface{}) {
	l.logger.Log(context.Background(), slog.Level(LevelTrace),
		fmt.Sprint(v...))
}

// Debug formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelDebug.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Debug(v ...interface{}) {
	l.logger.Log(context.Background(), slog.Level(LevelDebug),
		fmt.Sprint(v...))
}

// Info formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelInfo.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Info(v ...interface{}) {
	l.logger.Log(context.Background(), slog.Level(LevelInfo),
		fmt.Sprint(v...))
}

// Warn formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelWarn.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Warn(v ...interface{}) {
	l.logger.Log(context.Background(), slog.Level(LevelWarn),
		fmt.Sprint(v...))
}

// Error formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelError.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Error(v ...interface{}) {
	l.logger.Log(context.Background(), slog.Level(LevelError),
		fmt.Sprint(v...))
}

// Critical formats message using the default formats for its operands, prepends
// the prefix as necessary, and writes to log with LevelCritical.
//
// This is part of the Logger interface implementation.
func (l *sLogger) Critical(v ...interface{}) {
	l.logger.Log(context.Background(), slog.Level(LevelCritical),
		fmt.Sprint(v...))
}

// TraceS writes a structured log with the given message and key-value pair
// attributes with LevelTrace to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) TraceS(ctx context.Context, msg string, attrs ...any) {
	l.logger.Log(ctx, slog.Level(LevelTrace), msg, attrs...)
}

// DebugS writes a structured log with the given message and key-value pair
// attributes with LevelDebug to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) DebugS(ctx context.Context, msg string, attrs ...any) {
	l.logger.Log(context.Background(), slog.Level(LevelDebug), msg,
		mergeAttrs(ctx, attrs)...)
}

// InfoS writes a structured log with the given message and key-value pair
// attributes with LevelInfo to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) InfoS(ctx context.Context, msg string, attrs ...any) {
	l.logger.Log(ctx, slog.Level(LevelInfo), msg,
		mergeAttrs(ctx, attrs)...)
}

// WarnS writes a structured log with the given message and key-value pair
// attributes with LevelWarn to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) WarnS(ctx context.Context, msg string, err error,
	attrs ...any) {

	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.logger.Log(context.Background(), slog.Level(LevelWarn), msg,
		mergeAttrs(ctx, attrs)...)
}

// ErrorS writes a structured log with the given message and key-value pair
// attributes with LevelError to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) ErrorS(ctx context.Context, msg string, err error,
	attrs ...any) {

	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.logger.Log(context.Background(), slog.Level(LevelError), msg,
		mergeAttrs(ctx, attrs)...)
}

// CriticalS writes a structured log with the given message and key-value pair
// attributes with LevelCritical to the log.
//
// This is part of the Logger interface implementation.
func (l *sLogger) CriticalS(ctx context.Context, msg string, err error,
	attrs ...any) {
	if err != nil {
		attrs = append([]any{slog.String("err", err.Error())}, attrs...)
	}

	l.logger.Log(context.Background(), slog.Level(LevelCritical), msg,
		mergeAttrs(ctx, attrs)...)
}

var _ Logger = (*sLogger)(nil)
