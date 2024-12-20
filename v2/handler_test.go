package btclog

import (
	"bytes"
	"context"
	"errors"
	"github.com/btcsuite/btclog"
	"io"
	"testing"
	"time"
)

// TestDefaultHandler tests that the DefaultHandler's output looks as expected.
func TestDefaultHandler(t *testing.T) {
	t.Parallel()

	timeSource := func() time.Time {
		return time.Unix(100, 100)
	}

	tests := []struct {
		name               string
		handlerConstructor func(w io.Writer) Handler
		level              btclog.Level
		logFunc            func(log Logger)
		expectedLog        string
	}{
		{
			name: "Basic calls and levels",
			handlerConstructor: func(w io.Writer) Handler {
				return NewDefaultHandler(
					w, WithTimeSource(timeSource),
				)
			},
			level: LevelDebug,
			logFunc: func(log Logger) {
				log.Info("Test Basic Log")
				log.Debugf("Test basic log with %s", "format")
				log.Trace("Log should not appear due to level")
			},
			expectedLog: `1970-01-01 02:01:40.000 [INF]: Test Basic Log
1970-01-01 02:01:40.000 [DBG]: Test basic log with format
`,
		},
		{
			name: "Call site",
			handlerConstructor: func(w io.Writer) Handler {
				return NewDefaultHandler(
					w, WithNoTimestamp(),
					WithCallSiteSkipDepth(7),
					WithCallerFlags(Lshortfile),
				)
			},
			level: LevelInfo,
			logFunc: func(log Logger) {
				log.Info("Test Basic Log")
			},
			expectedLog: `[INF] handler_test.go:196: Test Basic Log
`,
		},
		{
			name: "Sub-system tag",
			handlerConstructor: func(w io.Writer) Handler {
				h := NewDefaultHandler(w, WithNoTimestamp())
				return h.SubSystem("SUBS")
			},
			level: LevelInfo,
			logFunc: func(log Logger) {
				log.Info("Test Basic Log")
			},
			expectedLog: `[INF] SUBS: Test Basic Log
`,
		},
		{
			name: "Test all levels",
			handlerConstructor: func(w io.Writer) Handler {
				return NewDefaultHandler(w, WithNoTimestamp())
			},
			level: LevelTrace,
			logFunc: func(log Logger) {
				log.Trace("Trace")
				log.Debug("Debug")
				log.Info("Info")
				log.Warn("Warn")
				log.Error("Error")
				log.Critical("Critical")
			},
			expectedLog: `[TRC]: Trace
[DBG]: Debug
[INF]: Info
[WRN]: Warn
[ERR]: Error
[CRT]: Critical
`,
		},
		{
			name: "Structured Logs",
			handlerConstructor: func(w io.Writer) Handler {
				return NewDefaultHandler(w, WithNoTimestamp())
			},
			level: LevelInfo,
			logFunc: func(log Logger) {
				ctx := context.Background()
				log.InfoS(ctx, "No attributes")
				log.InfoS(ctx, "Single word attribute", "key", "value")
				log.InfoS(ctx, "Multi word string value", "key with spaces", "value")
				log.InfoS(ctx, "Number attribute", "key", 5)
				log.InfoS(ctx, "Bad key", "key")

				type b struct {
					name    string
					age     int
					address *string
				}

				var c *b
				log.InfoS(ctx, "Nil pointer value", "key", c)

				c = &b{name: "Bob", age: 5}
				log.InfoS(ctx, "Struct values", "key", c)

				ctx = WithCtx(ctx, "request_id", 5, "user_name", "alice")
				log.InfoS(ctx, "Test context attributes", "key", "value")
			},
			expectedLog: `[INF]: No attributes
[INF]: Single word attribute key=value
[INF]: Multi word string value "key with spaces"=value
[INF]: Number attribute key=5
[INF]: Bad key !BADKEY=key
[INF]: Nil pointer value key=<nil>
[INF]: Struct values key="&{name:Bob age:5 address:<nil>}"
[INF]: Test context attributes request_id=5 user_name=alice key=value
`,
		},
		{
			name: "Error logs",
			handlerConstructor: func(w io.Writer) Handler {
				return NewDefaultHandler(w, WithNoTimestamp())
			},
			level: LevelInfo,
			logFunc: func(log Logger) {
				log.Error("Error string")
				log.Errorf("Error formatted string")

				ctx := context.Background()
				log.ErrorS(ctx, "Structured error log with nil error", nil)
				log.ErrorS(ctx, "Structured error with non-nil error", errors.New("oh no"))
				log.ErrorS(ctx, "Structured error with attributes", errors.New("oh no"), "key", "value")

				log.Warn("Warning string")
				log.Warnf("Warning formatted string")

				ctx = context.Background()
				log.WarnS(ctx, "Structured warning log with nil error", nil)
				log.WarnS(ctx, "Structured warning with non-nil error", errors.New("oh no"))
				log.WarnS(ctx, "Structured warning with attributes", errors.New("oh no"), "key", "value")

				log.Critical("Critical string")
				log.Criticalf("Critical formatted string")

				ctx = context.Background()
				log.CriticalS(ctx, "Structured critical log with nil error", nil)
				log.CriticalS(ctx, "Structured critical with non-nil error", errors.New("oh no"))
				log.CriticalS(ctx, "Structured critical with attributes", errors.New("oh no"), "key", "value")
			},
			expectedLog: `[ERR]: Error string
[ERR]: Error formatted string
[ERR]: Structured error log with nil error
[ERR]: Structured error with non-nil error err="oh no"
[ERR]: Structured error with attributes err="oh no" key=value
[WRN]: Warning string
[WRN]: Warning formatted string
[WRN]: Structured warning log with nil error
[WRN]: Structured warning with non-nil error err="oh no"
[WRN]: Structured warning with attributes err="oh no" key=value
[CRT]: Critical string
[CRT]: Critical formatted string
[CRT]: Structured critical log with nil error
[CRT]: Structured critical with non-nil error err="oh no"
[CRT]: Structured critical with attributes err="oh no" key=value
`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			handler := test.handlerConstructor(&buf)
			handler.SetLevel(test.level)

			if handler.Level() != test.level {
				t.Fatalf("Incorrect level. Expected %s, "+
					"got %s", test.level, handler.Level())
			}

			test.logFunc(NewSLogger(handler))

			if string(buf.Bytes()) != test.expectedLog {
				t.Fatalf("Log result mismatch. Expected "+
					"\n\"%s\", got \n\"%s\"",
					test.expectedLog, buf.Bytes())
			}
		})
	}
}
