// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package logger

import (
	"log/slog"
	"reflect"
	"testing"

	"github.com/urfave/cli/v3"
)

func TestLogLevel_String(t *testing.T) {
	var logl LogLevel
	tests := []struct {
		name  string
		ll    LogLevel
		wantS string
	}{
		{
			name:  "info",
			ll:    LogLevel(slog.LevelInfo),
			wantS: "INFO",
		},
		{
			name:  "error",
			ll:    LogLevel(slog.LevelError),
			wantS: "ERROR",
		},
		{
			name:  "warn",
			ll:    LogLevel(slog.LevelWarn),
			wantS: "WARN",
		},
		{
			name:  "debug",
			ll:    LogLevel(slog.LevelDebug),
			wantS: "DEBUG",
		},
		{
			name:  "trace",
			ll:    LogLevel(LevelTrace),
			wantS: "TRACE",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logl = tt.ll
			if gotS := logl.String(); gotS != tt.wantS {
				t.Errorf("LogLevel.String() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestLogLevel_Set(t *testing.T) {
	var logl LogLevel
	type args struct {
		ls string
	}
	tests := []struct {
		name    string
		ll      *LogLevel
		args    args
		wantErr bool
	}{
		{
			name: "info",
			ll:   &logl,
			args: args{
				ls: "info",
			},
			wantErr: false,
		},
		{
			name: "error",
			ll:   &logl,
			args: args{
				ls: "error",
			},
			wantErr: false,
		},
		{
			name: "warn",
			ll:   &logl,
			args: args{
				ls: "warn",
			},
			wantErr: false,
		},
		{
			name: "debug",
			ll:   &logl,
			args: args{
				ls: "debug",
			},
			wantErr: false,
		},
		{
			name: "trace",
			ll:   &logl,
			args: args{
				ls: "trace",
			},
			wantErr: false,
		},
		{
			name: "fail",
			ll:   &logl,
			args: args{
				ls: "fail",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ll.Set(tt.args.ls); (err != nil) != tt.wantErr {
				t.Errorf("LogLevel.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLogLevel_Type(t *testing.T) {
	var logl LogLevel
	tests := []struct {
		name string
		ll   *LogLevel
		want string
	}{
		{
			name: "ok",
			ll:   &logl,
			want: "LogLevel",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.ll.Type(); got != tt.want {
				t.Errorf("LogLevel.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogLevel_UnmarshalJSON(t *testing.T) {
	type args struct {
		jason []byte
	}
	var ll LogLevel
	tests := []struct {
		name    string
		ll      *LogLevel
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			ll:   &ll,
			args: args{
				jason: []byte("warn"),
			},
			wantErr: false,
		},
		{
			name: "error",
			ll:   &ll,
			args: args{
				jason: []byte("freddy"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.ll.UnmarshalJSON(tt.args.jason); (err != nil) != tt.wantErr {
				t.Errorf("LogLevel.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_logLevelValue_Create(t *testing.T) {
	var (
		ll  LogLevel
		llv = logLevelValue{
			destination: &ll,
		}
	)
	type args struct {
		val LogLevel
		p   *LogLevel
		in2 cli.NoConfig
	}
	tests := []struct {
		name string
		args args
		want cli.Value
	}{
		{
			name: "ok",
			args: args{
				val: LogLevel(LevelTrace),
				p:   &ll,
				in2: cli.NoConfig{},
			},
			want: llv,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logLevelValue{}
			if got := l.Create(tt.args.val, tt.args.p, tt.args.in2); got.String() != tt.want.String() {
				t.Errorf("logLevelValue.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logLevelValue_Get(t *testing.T) {
	type fields struct {
		destination *LogLevel
	}
	var ll = LogLevel(slog.LevelInfo)
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "ok",
			fields: fields{
				destination: &ll,
			},
			want: LogLevel(slog.LevelInfo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logLevelValue{
				destination: tt.fields.destination,
			}
			if got := l.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("logLevelValue.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logLevelValue_Set(t *testing.T) {
	type fields struct {
		destination *LogLevel
	}
	type args struct {
		s string
	}
	var ll LogLevel
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				destination: &ll,
			},
			args: args{
				s: "INFO",
			},
			wantErr: false,
		},
		{
			name: "error",
			fields: fields{
				destination: &ll,
			},
			args: args{
				s: "fred",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logLevelValue{
				destination: tt.fields.destination,
			}
			if err := l.Set(tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("logLevelValue.Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_logLevelValue_String(t *testing.T) {
	type fields struct {
		destination *LogLevel
	}
	var ll = LogLevel(slog.LevelError)
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				destination: &ll,
			},
			want: "ERROR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logLevelValue{
				destination: tt.fields.destination,
			}
			if got := l.String(); got != tt.want {
				t.Errorf("logLevelValue.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_logLevelValue_ToString(t *testing.T) {
	type fields struct {
		destination *LogLevel
	}
	type args struct {
		lev LogLevel
	}
	var ll LogLevel
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "ok",
			fields: fields{
				destination: &ll,
			},
			args: args{
				lev: LogLevel(slog.LevelError),
			},
			want: "ERROR",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := logLevelValue{
				destination: tt.fields.destination,
			}
			if got := l.ToString(tt.args.lev); got != tt.want {
				t.Errorf("logLevelValue.ToString() = %v, want %v", got, tt.want)
			}
		})
	}
}
