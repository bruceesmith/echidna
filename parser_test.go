// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package echidna

import (
	"reflect"
	"testing"
)

func Test_camelToFlag(t *testing.T) {
	type args struct {
		s           string
		flagDivider string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "one",
			args: args{
				s:           "ThisIsACamel",
				flagDivider: "-",
			},
			want: "this-is-a-camel",
		},
		{
			name: "two",
			args: args{
				s:           "notacamel",
				flagDivider: "-",
			},
			want: "notacamel",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := camelToFlag(tt.args.s, tt.args.flagDivider); got != tt.want {
				t.Errorf("camelToFlag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseField(t *testing.T) {
	type config2 struct {
		I int `flag:"-"`
		J int `flag:""`
		K int `flag:"kk"`
		L int `flag:"el elel"`
		M int `flag:"~em eminem"`
	}
	var (
		fld, fld1, fld2, fld3, fld4, fld5 reflect.StructField
	)
	type args struct {
		field  reflect.StructField
		prefix string
	}

	fld, _ = reflect.TypeFor[config]().FieldByName("I")
	fld1, _ = reflect.TypeFor[config2]().FieldByName("I")
	fld2, _ = reflect.TypeFor[config2]().FieldByName("J")
	fld3, _ = reflect.TypeFor[config2]().FieldByName("K")
	fld4, _ = reflect.TypeFor[config2]().FieldByName("L")
	fld5, _ = reflect.TypeFor[config2]().FieldByName("M")
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "no-tag",
			args: args{
				field: fld,
			},
			want: "i",
		},
		{
			name: "hyphen-tag",
			args: args{
				field: fld1,
			},
			want: "",
		},
		{
			name: "empty-tag",
			args: args{
				field: fld2,
			},
			want: "j",
		},
		{
			name: "alternate-name",
			args: args{
				field: fld3,
			},
			want: "kk",
		},
		{
			name: "multiple-names",
			args: args{
				field: fld4,
			},
			want: "el",
		},
		{
			name: "remove prefix",
			args: args{
				field:  fld5,
				prefix: "something",
			},
			want: "em",
		},
	}
	for _, tt := range tests {
		opts = options{
			divider: "-",
			tag:     "flag",
			flatten: false,
			prefix:  "",
		}
		t.Run(tt.name, func(t *testing.T) {
			if got := parseField(tt.args.field, tt.args.prefix); got != tt.want {
				t.Errorf("parseField() = %v, want %v", got, tt.want)
			}
		})
	}
}
