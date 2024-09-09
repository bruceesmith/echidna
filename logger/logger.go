// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/bruceesmith/echidna/set"
)

const (
	// LevelTrace can be set to enable tracing
	LevelTrace slog.Level = -10
)

var (
	level    slog.LevelVar
	traceIds *set.Set[string]
	trace    *slog.Logger
)

func init() {
	level.Set(slog.LevelInfo)
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{
					Level: &level,
				},
			),
		),
	)
	traceIds = set.New[string]()
	trace = slog.New(
		slog.NewJSONHandler(
			os.Stdout,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     &level,
			},
		),
	)
}

// Debug emits a debug log
func Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

// Error emits an error log
func Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

// Info emits an info log
func Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func Level() string {
	return level.String()
}

// RedirectStandard changes the destination for normal (non-trace) logs
func RedirectStandard(w io.Writer) {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				w,
				&slog.HandlerOptions{
					Level: &level,
				},
			),
		),
	)
}

// RedirectTrace changes the destination for normal (non-trace) logs
func RedirectTrace(w io.Writer) {
	trace = slog.New(
		slog.NewJSONHandler(
			w,
			&slog.HandlerOptions{
				AddSource: true,
				Level:     &level,
			},
		),
	)
}

// SetLevel sets the default level of logging
func SetLevel(l slog.Level) {
	level.Set(l)
}

// SetTraceIds registers identifiers for future tracing
func SetTraceIds(ids ...string) {
	for _, id := range ids {
		traceIds.Add(strings.ToLower(id))
	}
}

// Trace emits one JSON-formatted log entry if trace level logging is enabled
func Trace(msg string, args ...any) {
	ctx := context.Background()
	trace.Log(ctx, LevelTrace, msg, args...)
}

// TraceID emits one JSON-formatted log entry if tracing is enabled for the requested ID
func TraceID(id string, msg string, args ...any) {
	if traceIds.Contains(strings.ToLower(id)) || traceIds.Contains("all") {
		ctx := context.Background()
		trace.Log(ctx, LevelTrace, msg, args...)
	}
}

// Warn emits a warning log
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}
