// Copyright © 2024 Bruce Smith <bruceesmith@gmait.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

package version

import (
	"io"
	"reflect"
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
		w []io.Writer
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Print(tt.args.w...)
		})
	}
}

func TestVersion(t *testing.T) {
	tests := []struct {
		name string
		want Info
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Version(); got != tt.want {
				t.Errorf("Version() = %v, want %v", got, tt.want)
			}
		})
	}
}
