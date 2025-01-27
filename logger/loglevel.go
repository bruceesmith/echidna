// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package logger

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/urfave/cli/v3"
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
		err = fmt.Errorf("invalid log level %v", ls)
	}
	return
}

// Type is a conveniene method for pflag.Value
func (ll *LogLevel) Type() string {
	return "LogLevel"
}

// UnmarshalJSON is a convenience method for Kong
func (ll *LogLevel) UnmarshalJSON(jason []byte) (err error) {
	err = ll.Set(string(jason))
	if err != nil {
		return fmt.Errorf("cannot unmarshal %s, %w", string(jason), err)
	}
	return
}

// LogLevelFlag is useful for using a LogLevel as a command-line flag in CLI applications
type LogLevelFlag = cli.FlagBase[LogLevel, cli.NoConfig, logLevelValue]

// logLevelValue supports command-line LogLevel arguments
type logLevelValue struct {
	destination *LogLevel
}

// Create returns a value which implements the golang flag.Value and flag.Getter interfaces
func (l logLevelValue) Create(val LogLevel, p *LogLevel, _ cli.NoConfig) cli.Value {
	*p = val
	return &logLevelValue{destination: p}
}

// Get fetches the LogLevel value
func (l logLevelValue) Get() any {
	return *l.destination
}

// Set stores a string into a LogLevelFlag
func (l logLevelValue) Set(s string) error {
	return l.destination.Set(s)
}

// String returns a string representation of a LogLevel
func (l logLevelValue) String() string {
	return l.destination.String()
}

// ToString returns a string representation of a LogLevel
func (l logLevelValue) ToString(lev LogLevel) string {
	return lev.String()
}
