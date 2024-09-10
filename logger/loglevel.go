// Copyright © 2024 Bruce Smith <bruceesmith@gmait.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

package logger

import (
	"log/slog"
	"strings"
)

// LogLevel is the level of logging
type LogLevel int

// String is a convenience method for pflag.Value
func (ll *LogLevel) String() (s string) {
	switch *ll {
	case LogLevel(slog.LevelInfo):
		s = "INFO"
	case LogLevel(slog.LevelError):
		s = "ERROR"
	case LogLevel(slog.LevelWarn):
		s = "WARN"
	case LogLevel(slog.LevelDebug):
		s = "DEBUG"
	case LogLevel(LevelTrace):
		s = "TRACE"
	}
	return
}

// Set is a convenience method for pflag.Value
func (ll *LogLevel) Set(ls string) (err error) {
	switch strings.ToLower(ls) {
	case "info":
		*ll = LogLevel(slog.LevelInfo)
	case "error":
		*ll = LogLevel(slog.LevelError)
	case "warn":
		*ll = LogLevel(slog.LevelWarn)
	case "debug":
		*ll = LogLevel(slog.LevelDebug)
	case "trace":
		*ll = LogLevel(LevelTrace)
	default:
		*ll = LogLevel(slog.LevelInfo)
	}
	return
}

// Type is a conveniene method for pflag.Value
func (ll *LogLevel) Type() string {
	return "LogLevel"
}
