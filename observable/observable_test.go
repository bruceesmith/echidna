// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package observable

import (
	"testing"
)

func TestObserve(t *testing.T) {
	type args struct {
		name  string
		cb    func(int)
		prev  func(name string, cb func(any)) error
		pname string
		pcb   func(int)
		ncb   func(string)
	}
	var (
		f1 = func(i int) {}
		f2 = func(i int) {}
		f3 = func(s string) {}
	)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok-first-registration",
			args: args{
				name:  "fred",
				cb:    f1,
				pname: "",
			},
			wantErr: false,
		},
		{
			name: "ok-second-registration",
			args: args{
				name:  "fred",
				cb:    f2,
				pname: "fred",
				pcb:   f1,
			},
			wantErr: false,
		},
		{
			name: "error-second-registration",
			args: args{
				name:  "fred",
				ncb:   f3,
				pname: "fred",
				pcb:   f1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.pname != "" {
				_ = Observe(tt.args.pname, tt.args.pcb)
			}
			if tt.args.cb != nil {
				if err := Observe(tt.args.name, tt.args.cb); (err != nil) != tt.wantErr {
					t.Errorf("Observe() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else if tt.args.ncb != nil {
				if err := Observe(tt.args.name, tt.args.ncb); (err != nil) != tt.wantErr {
					t.Errorf("Observe() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestSet(t *testing.T) {
	type args struct {
		name   string
		value  int
		svalue string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				name:  "fred",
				value: 1,
			},
			wantErr: false,
		},
		{
			name: "error-wrong-type",
			args: args{
				name:   "fred",
				svalue: "bill",
			},
			wantErr: true,
		},
		{
			name: "error-not-registered",
			args: args{
				name:  "mary",
				value: 1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Observe("fred", func(i int) {})
			if err != nil {
				t.Errorf("Set() error = %v in setup", err)
			}
			if len(tt.args.svalue) == 0 {
				if err = Set(tt.args.name, tt.args.value); (err != nil) != tt.wantErr {
					t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if err = Set(tt.args.name, tt.args.svalue); (err != nil) != tt.wantErr {
					t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}
