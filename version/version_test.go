// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package version

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"regexp"
	"testing"
)

func Test_makeVersion(t *testing.T) {
	tests := []struct {
		name   string
		wantVi Info
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotVi := makeVersion(); !reflect.DeepEqual(gotVi, tt.wantVi) {
				t.Errorf("makeVersion() = %v, want %v", gotVi, tt.wantVi)
			}
		})
	}
}

func TestPrint(t *testing.T) {
	type args struct {
		jason bool
		w     []io.Writer
	}
	var buf = bytes.NewBufferString("")
	tests := []struct {
		name string
		args args
	}{
		{
			name: "fail",
			args: args{},
		},
		{
			name: "plain text",
			args: args{
				jason: false,
				w:     []io.Writer{buf},
			},
		},
		{
			name: "json",
			args: args{
				jason: true,
				w:     []io.Writer{buf},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.args.w) == 0 {
				return
			}
			buf.Reset()
			Print(tt.args.jason, tt.args.w...)
			if !tt.args.jason {
				matched, err := regexp.MatchString(`.*\nCompiled.+\nBuilt.+`, buf.String())
				if !matched || err != nil {
					t.Errorf("Print(): matched %v error %v", matched, err)
				}
			} else {
				var inf Info
				err := json.Unmarshal(buf.Bytes(), &inf)
				if err != nil {
					t.Errorf("Print(): %v", err)
				}
				if len(inf.Error) > 0 {
					t.Errorf("Print(): %v", inf.Error)
				}
			}
		})
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		want Info
	}{
		{
			name: "ok",
			want: makeVersion(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}
