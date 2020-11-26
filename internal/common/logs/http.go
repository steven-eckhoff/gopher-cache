// Package logs implements a library that provides common logging utilities.
package logs

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/sirupsen/logrus"
)

// NewStructuredLogger creates a new go-chi middleware using a newly created StructuredLogger.
func NewStructuredLogger(logger *logrus.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{logger})
}

// based on example from chi: https://github.com/go-chi/chi/blob/master/_examples/logging/main.go
type StructuredLogger struct {
	Logger *logrus.Logger
}

// NewLogEntry creates a log entry for a given request.
func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	entry := &StructuredLoggerEntry{Logger: logrus.NewEntry(l.Logger)}
	logFields := logrus.Fields{}

	if reqID := middleware.GetReqID(r.Context()); reqID != "" {
		logFields["req_id"] = reqID
	}

	logFields["http_method"] = r.Method

	logFields["remote_addr"] = r.RemoteAddr
	logFields["uri"] = r.RequestURI

	entry.Logger = entry.Logger.WithFields(logFields)

	entry.Logger.Info("Request started")

	return entry
}

// StructuredLoggerEntry represents one log entry.
type StructuredLoggerEntry struct {
	Logger logrus.FieldLogger
}

func (l *StructuredLoggerEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"resp_status": status, "resp_bytes_length": bytes,
		"resp_elapsed": elapsed.Round(time.Millisecond / 100).String(),
	})

	l.Logger.Info("Request completed	")
}

func (l *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	l.Logger = l.Logger.WithFields(logrus.Fields{
		"stack": string(stack),
		"panic": fmt.Sprintf("%+v", v),
	})
}

// GetLogEntry returns the log entry associated with the given request.
func GetLogEntry(r *http.Request) logrus.FieldLogger {
	entry := middleware.GetLogEntry(r).(*StructuredLoggerEntry)
	return entry.Logger
}
