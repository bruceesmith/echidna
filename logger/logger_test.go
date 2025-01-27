// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package logger

import (
	"bytes"
	"log/slog"
	"regexp"
	"strings"
	"testing"

	"github.com/bruceesmith/echidna/set"
)

func TestDebug(t *testing.T) {
	type args struct {
		msg  string
		args []any
	}
	tests := []struct {
		name   string
		args   args
		level  slog.Level
		wantRe string
	}{
		{
			name: "debug-ok",
			args: args{
				msg:  "debug",
				args: []any{"one", 1},
			},
			level:  slog.LevelDebug,
			wantRe: "^time=.+ level=DEBUG msg=debug one=1",
		},
		{
			name: "debug-below-level",
			args: args{
				msg:  "debug",
				args: []any{"one", 1},
			},
			level:  slog.LevelInfo,
			wantRe: "^$",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			slog.SetDefault(
				slog.New(
					slog.NewTextHandler(
						w,
						&slog.HandlerOptions{
							Level: &tt.level,
						},
					),
				),
			)
			Debug(tt.args.msg, tt.args.args...)
			ok, err := regexp.MatchString(tt.wantRe, w.String())
			if !ok {
				t.Errorf("Debug() got %s want %s error %s", w.String(), tt.wantRe, err)
			}
		})
	}
}

func TestError(t *testing.T) {
	type args struct {
		msg  string
		args []any
	}
	tests := []struct {
		name   string
		args   args
		level  slog.Level
		wantRe string
	}{
		{
			name: "error-ok",
			args: args{
				msg:  "error",
				args: []any{"one", 1},
			},
			level:  slog.LevelError,
			wantRe: "^time=.+ level=ERROR msg=error one=1",
		},
		{
			name: "error-below-level",
			args: args{
				msg:  "error",
				args: []any{"one", 1},
			},
			level:  slog.LevelError + 1,
			wantRe: "^$",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			slog.SetDefault(
				slog.New(
					slog.NewTextHandler(
						w,
						&slog.HandlerOptions{
							Level: &tt.level,
						},
					),
				),
			)
			Error(tt.args.msg, tt.args.args...)
			ok, err := regexp.MatchString(tt.wantRe, w.String())
			if !ok {
				t.Errorf("Error() got %s want %s error %s", w.String(), tt.wantRe, err)
			}
		})
	}
}

func TestInfo(t *testing.T) {
	type args struct {
		msg  string
		args []any
	}
	tests := []struct {
		name   string
		args   args
		level  slog.Level
		wantRe string
	}{
		{
			name: "info-ok",
			args: args{
				msg:  "info",
				args: []any{"one", 1},
			},
			level:  slog.LevelInfo,
			wantRe: "^time=.+ level=INFO msg=info one=1",
		},
		{
			name: "info-below-level",
			args: args{
				msg:  "info",
				args: []any{"one", 1},
			},
			level:  slog.LevelWarn,
			wantRe: "^$",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			slog.SetDefault(
				slog.New(
					slog.NewTextHandler(
						w,
						&slog.HandlerOptions{
							Level: &tt.level,
						},
					),
				),
			)
			Info(tt.args.msg, tt.args.args...)
			ok, err := regexp.MatchString(tt.wantRe, w.String())
			if !ok {
				t.Errorf("Info() got %s want %s error %s", w.String(), tt.wantRe, err)
			}
		})
	}
}

func Test_jsonHandler(t *testing.T) {
	type args struct {
		trace bool
	}
	tests := []struct {
		name string
		args args
	}{{
		name: "trace",
		args: args{
			trace: true,
		},
	},

		{
			name: "no-trace",
			args: args{
				trace: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			if got := jsonHandler(w, tt.args.trace); got == nil {
				t.Error("jsonHandler() returned nil")
			}

		})
	}
}

func TestLevel(t *testing.T) {
	tests := []struct {
		name string
		lev  slog.Level
		want string
	}{
		{
			name: "info",
			lev:  slog.LevelInfo,
			want: "INFO",
		},
		{
			name: "trace",
			lev:  LevelTrace,
			want: "DEBUG-6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.lev)
			if got := Level(); got != tt.want {
				t.Errorf("Level() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRedirectStandard(t *testing.T) {
	tests := []struct {
		name   string
		format Format
	}{
		{
			name:   "json",
			format: JSON,
		},
		{
			name:   "text",
			format: Text,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format = tt.format
			w := &bytes.Buffer{}
			before := slog.Default()
			RedirectStandard(w)
			after := slog.Default()
			if before == after {
				t.Errorf("RedirectStandard() before = %v, after = %v", before, after)
			}
		})
	}
}

func TestRedirectTrace(t *testing.T) {
	tests := []struct {
		name   string
		format Format
	}{
		{
			name:   "json",
			format: JSON,
		},
		{
			name:   "text",
			format: Text,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			format = tt.format
			w := &bytes.Buffer{}
			before := *trace
			RedirectTrace(w)
			after := *trace
			if before == after {
				t.Errorf("RedirectTrace() before = %v, after = %v", before, after)
			}
		})
	}
}

func TestSetFormat(t *testing.T) {
	type args struct {
		f Format
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "json",
			args: args{
				f: JSON,
			},
		},
		{
			name: "text",
			args: args{
				f: Text,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetFormat(tt.args.f)
		})
	}
}

func TestSetLevel(t *testing.T) {
	type args struct {
		l slog.Level
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "debug",
			args: args{
				l: slog.LevelDebug,
			},
		},
		{
			name: "trace",
			args: args{
				l: LevelTrace,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.args.l)
			if level.Level() != tt.args.l {
				t.Errorf("SetLevel() = %v, want %v", level.Level(), tt.args.l)
			}
		})
	}
}

func TestSetTraceIds(t *testing.T) {
	type args struct {
		ids []string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple",
			args: args{
				ids: []string{"one", "two"},
			},
		},
		{
			name: "mixed case",
			args: args{
				ids: []string{"oNe", "Two"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetTraceIds(tt.args.ids...)
			for _, id := range tt.args.ids {
				if !traceIds.Contains(strings.ToLower(id)) {
					t.Errorf("SetTraceIds %s is not in traceIds", id)
				}
			}
		})
	}
}

func TestTrace(t *testing.T) {
	type args struct {
		msg  string
		args []any
	}
	tests := []struct {
		name   string
		args   args
		level  slog.Level
		wantRe string
	}{
		{
			name: "trace-ok",
			args: args{
				msg:  "trace",
				args: []any{"one", 1},
			},
			level:  LevelTrace,
			wantRe: `^{"time":".+,"level":"DEBUG-6","msg":"trace","one":1}`,
		},
		{
			name: "trace-below-level",
			args: args{
				msg:  "trace",
				args: []any{"one", 1},
			},
			level:  slog.LevelInfo,
			wantRe: "^$",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.level)
			w := &bytes.Buffer{}
			trace =
				slog.New(
					slog.NewJSONHandler(
						w,
						&slog.HandlerOptions{
							Level: &tt.level,
						},
					),
				)
			Trace(tt.args.msg, tt.args.args...)
			s := w.String()
			ok, err := regexp.MatchString(tt.wantRe, s)
			if !ok {
				t.Errorf("Trace() got %s want %s error %v", s, tt.wantRe, err)
			}
		})
	}
}

func TestTraceID(t *testing.T) {
	type args struct {
		msg  string
		args []any
		ids  *set.Set[string]
	}
	tests := []struct {
		name   string
		args   args
		level  slog.Level
		id     string
		wantRe string
	}{
		{
			name: "trace-bad-id",
			args: args{
				msg:  "trace",
				args: []any{"one", 1},
				ids:  set.New[string]("m1", "m2"),
			},
			level:  LevelTrace,
			id:     "m3",
			wantRe: `^$`,
		},
		{
			name: "trace-ok",
			args: args{
				msg:  "trace",
				args: []any{"one", 1},
				ids:  set.New[string]("m1", "m2"),
			},
			level:  LevelTrace,
			id:     "m1",
			wantRe: `^{"time":".+,"level":"DEBUG-6","msg":"trace","one":1}`,
		},
		{
			name: "trace-below-level",
			args: args{
				msg:  "trace",
				args: []any{"one", 1},
				ids:  set.New[string]("m1", "m2"),
			},
			level:  slog.LevelInfo,
			id:     "m1",
			wantRe: "^$",
		},
		{
			name: "trace-all",
			args: args{
				msg:  "trace",
				args: []any{"one", 1},
				ids:  set.New[string]("all"),
			},
			level:  LevelTrace,
			id:     "m1",
			wantRe: `^{"time":".+,"level":"DEBUG-6","msg":"trace","one":1}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			SetLevel(tt.level)
			trace =
				slog.New(
					slog.NewJSONHandler(
						w,
						&slog.HandlerOptions{
							Level: &tt.level,
						},
					),
				)
			traceIds = tt.args.ids
			TraceID(tt.id, tt.args.msg, tt.args.args...)
			s := w.String()
			ok, err := regexp.MatchString(tt.wantRe, s)
			if !ok {
				t.Errorf("Trace() got %s want %s error %s", s, tt.wantRe, err)
			}
		})
	}
}

func TestWarn(t *testing.T) {
	type args struct {
		msg  string
		args []any
	}
	tests := []struct {
		name   string
		args   args
		level  slog.Level
		wantRe string
	}{
		{
			name: "warn-ok",
			args: args{
				msg:  "warn",
				args: []any{"one", 1},
			},
			level:  slog.LevelInfo,
			wantRe: "^time=.+ level=WARN msg=warn one=1",
		},
		{
			name: "warn-below-level",
			args: args{
				msg:  "warn",
				args: []any{"one", 1},
			},
			level:  slog.LevelError,
			wantRe: "^$",
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			slog.SetDefault(
				slog.New(
					slog.NewTextHandler(
						w,
						&slog.HandlerOptions{
							Level: &tt.level,
						},
					),
				),
			)
			Warn(tt.args.msg, tt.args.args...)
			ok, err := regexp.MatchString(tt.wantRe, w.String())
			if !ok {
				t.Errorf("Warn() got %s want %s error %s", w.String(), tt.wantRe, err)
			}
		})
	}
}
