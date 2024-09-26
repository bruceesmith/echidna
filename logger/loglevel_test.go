// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

package logger

import (
	"log/slog"
	"testing"
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
