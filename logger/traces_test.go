// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package logger

import (
	"reflect"
	"testing"
)

func TestTraces_String(t *testing.T) {
	var trc Traces
	tests := []struct {
		name  string
		tr    []string
		wantS string
	}{
		{
			name:  "one",
			tr:    []string{"one", "two"},
			wantS: "[one two]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			trc = tt.tr
			if gotS := trc.String(); gotS != tt.wantS {
				t.Errorf("Traces.String() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestTraces_Set(t *testing.T) {
	var trc Traces
	type args struct {
		ts string
	}
	tests := []struct {
		name    string
		want    Traces
		args    args
		wantErr bool
	}{
		{
			name: "error",
			want: Traces{"area1", "area2"},
			args: args{
				ts: "area1,area2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := trc.Set(tt.args.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("Traces.Set() error = %v, wantErr %v", err, tt.wantErr)
			} else if !reflect.DeepEqual(trc, tt.want) {
				t.Errorf("Traces.Set() got = %v, want %v", trc, tt.want)
			}
		})
	}
}

func TestTraces_Type(t *testing.T) {
	var trc Traces
	tests := []struct {
		name string
		tr   *Traces
		want string
	}{
		{
			name: "ok",
			tr:   &trc,
			want: "Traces",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.tr.Type(); got != tt.want {
				t.Errorf("Traces.Type() = %v, want %v", got, tt.want)
			}
		})
	}
}
