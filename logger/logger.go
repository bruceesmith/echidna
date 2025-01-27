// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

/*
Package logger supports logging and tracing based on the standard library package slog.

Debug, Error, Info and Warn operate like their package slog equivalents, with the level of logging modifiable
using SetLevel. A custom logging level (LevelTrace) can be supplied to SetLevel to enable tracing. Tracing can
be unconditional when calling Trace, or only enabled for pre-defined identifiers when calling TraceID. Identifiers
for TraceID are registered by calling SetTraceIDs. By default, all debug, error, info and warn messages go to
Stdout, and traces go to Stderr; these destinations can be changed by calling RedirectNormal and RedirectTrace
respectively
*/
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/bruceesmith/echidna/set"
)

// Format determines the format of each log entry
type Format string

const (
	// LevelTrace can be set to enable tracing
	LevelTrace slog.Level = -10
	// Text format
	Text Format = "text"
	// JSON format
	JSON Format = "json"
)

var (
	format       Format
	level        slog.LevelVar
	traceIds     *set.Set[string]
	trace        *slog.Logger
	normalWriter io.Writer
	traceWriter  io.Writer
)

func init() {
	format = Text
	level.Set(slog.LevelInfo)
	normalWriter, traceWriter = os.Stdout, os.Stdout
	slog.SetDefault(
		slog.New(
			textHandler(normalWriter),
		),
	)
	traceIds = set.New[string]()
	trace = slog.New(
		textHandler(traceWriter),
	)
}

func jsonHandler(w io.Writer, trace bool) slog.Handler {
	return slog.NewJSONHandler(
		w,
		&slog.HandlerOptions{
			AddSource: trace,
			Level:     &level,
		},
	)
}
func textHandler(w io.Writer) slog.Handler {
	return slog.NewTextHandler(
		w,
		&slog.HandlerOptions{
			Level: &level,
		},
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
	return level.Level().String()
}

// RedirectStandard changes the destination for normal (non-trace) logs
func RedirectStandard(w io.Writer) {
	normalWriter = w
	switch format {
	case JSON:
		slog.SetDefault(slog.New(jsonHandler(w, false)))
	case Text:
		slog.SetDefault(slog.New(textHandler(w)))
	}
}

// RedirectTrace changes the destination for normal (non-trace) logs
func RedirectTrace(w io.Writer) {
	traceWriter = w
	switch format {
	case JSON:
		trace = slog.New(jsonHandler(w, true))
	case Text:
		trace = slog.New(textHandler(w))
	}
}

// SetFormat changes the format of log entries
func SetFormat(f Format) {
	switch strings.ToLower(string(f)) {
	case "json":
		format = JSON
		slog.SetDefault(slog.New(jsonHandler(normalWriter, false)))
		trace = slog.New(jsonHandler(traceWriter, true))
	case "text":
		format = Text
		slog.SetDefault(slog.New(textHandler(normalWriter)))
		trace = slog.New(textHandler(traceWriter))
	}
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
	if level.Level() == LevelTrace {
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		r := slog.NewRecord(time.Now(), LevelTrace, msg, pcs[0])
		r.Add(args...)
		_ = trace.Handler().Handle(context.Background(), r)
	}
}

// TraceID emits one JSON-formatted log entry if tracing is enabled for the requested ID
func TraceID(id string, msg string, args ...any) {
	if level.Level() == LevelTrace && (traceIds.Contains(strings.ToLower(id)) || traceIds.Contains("all")) {
		var pcs [1]uintptr
		runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
		r := slog.NewRecord(time.Now(), LevelTrace, msg, pcs[0])
		r.Add(args...)
		_ = trace.Handler().Handle(context.Background(), r)
	}
}

// Warn emits a warning log
func Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}
