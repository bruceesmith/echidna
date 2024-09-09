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

func TestVersion(t *testing.T) {
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
			Version(tt.args.w...)
		})
	}
}

func TestJSON(t *testing.T) {
	tests := []struct {
		name      string
		wantJason string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotJason := JSON(); gotJason != tt.wantJason {
				t.Errorf("JSON() = %v, want %v", gotJason, tt.wantJason)
			}
		})
	}
}
