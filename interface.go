// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package btclog

import "context"

// Logger is an interface which describes a level-based logger.  A default
// implementation of Logger is implemented by this package and can be created
// by calling (*Backend).Logger.
type Logger interface {
	// Tracef formats message according to format specifier and writes to
	// to log with LevelTrace.
	Tracef(format string, params ...interface{})

	// Debugf formats message according to format specifier and writes to
	// log with LevelDebug.
	Debugf(format string, params ...interface{})

	// Infof formats message according to format specifier and writes to
	// log with LevelInfo.
	Infof(format string, params ...interface{})

	// Warnf formats message according to format specifier and writes to
	// to log with LevelWarn.
	Warnf(format string, params ...interface{})

	// Errorf formats message according to format specifier and writes to
	// to log with LevelError.
	Errorf(format string, params ...interface{})

	// Criticalf formats message according to format specifier and writes to
	// log with LevelCritical.
	Criticalf(format string, params ...interface{})

	// Trace formats message using the default formats for its operands
	// and writes to log with LevelTrace.
	Trace(v ...interface{})

	// Debug formats message using the default formats for its operands
	// and writes to log with LevelDebug.
	Debug(v ...interface{})

	// Info formats message using the default formats for its operands
	// and writes to log with LevelInfo.
	Info(v ...interface{})

	// Warn formats message using the default formats for its operands
	// and writes to log with LevelWarn.
	Warn(v ...interface{})

	// Error formats message using the default formats for its operands
	// and writes to log with LevelError.
	Error(v ...interface{})

	// Critical formats message using the default formats for its operands
	// and writes to log with LevelCritical.
	Critical(v ...interface{})

	// TraceS writes a structured log with the given message and key-value
	// pair attributes with LevelTrace to the log.
	TraceS(ctx context.Context, msg string, attrs ...any)

	// DebugS writes a structured log with the given message and key-value
	// pair attributes with LevelDebug to the log.
	DebugS(ctx context.Context, msg string, attrs ...any)

	// InfoS writes a structured log with the given message and key-value
	// pair attributes with LevelInfo to the log.
	InfoS(ctx context.Context, msg string, attrs ...any)

	// WarnS writes a structured log with the given message and key-value
	// pair attributes with LevelWarn to the log.
	WarnS(ctx context.Context, msg string, err error, attrs ...any)

	// ErrorS writes a structured log with the given message and key-value
	// pair attributes with LevelError to the log.
	ErrorS(ctx context.Context, msg string, err error, attrs ...any)

	// CriticalS writes a structured log with the given message and
	// key-value pair attributes with LevelCritical to the log.
	CriticalS(ctx context.Context, msg string, err error, attrs ...any)

	// Level returns the current logging level.
	Level() Level

	// SetLevel changes the logging level to the passed level.
	SetLevel(level Level)
}
