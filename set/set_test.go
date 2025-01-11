// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package set

import (
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	type E = int
	type args struct {
		vals []E
	}
	tests := []struct {
		name    string
		args    args
		want    *Set[E]
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				vals: []E{1, 2, 3},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					1: {},
					2: {},
					3: {},
				},
			},
			wantErr: false,
		},
		{
			name: "bad",
			args: args{
				vals: []E{1, 2, 3},
			},
			want:    &Set[E]{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.vals...)
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_set_Add(t *testing.T) {
	type E = int
	type args struct {
		in   []E
		vals []E
	}
	tests := []struct {
		name    string
		args    args
		want    *Set[E]
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				in:   []E{1, 2},
				vals: []E{3},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					1: {},
					2: {},
					3: {},
				},
			},
			wantErr: false,
		},
		{
			name: "bad",
			args: args{
				in:   []E{1, 2},
				vals: []E{3},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					1: {},
					2: {},
					4: {},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.args.in...)
			s.Add(tt.args.vals...)
			if !reflect.DeepEqual(s, tt.want) && !tt.wantErr {
				t.Errorf("Add() = %v, want %v", s, tt.want)
			}
		})
	}
}

func Test_set_Clear(t *testing.T) {
	type E = int
	type args struct {
		in []E
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				in: []E{1, 2},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.args.in...)
			s.Clear()
			if s.Size() != 0 {
				t.Errorf("Clear() = %v want empty", s)
			}
		})
	}
}

func Test_set_Contains(t *testing.T) {
	type E = int
	type args struct {
		in  []E
		val E
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				in:  []E{1, 2},
				val: 2,
			},

			want: true,
		},
		{
			name: "bad",
			args: args{
				in:  []E{1, 2},
				val: 4,
			},

			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.args.in...)
			if s.Contains(tt.args.val) != tt.want {
				t.Errorf("Contains() = %v want %v", s, tt.want)
			}
		})
	}
}

func Test_set_Delete(t *testing.T) {
	type E = int
	type args struct {
		in   []E
		vals []E
	}
	tests := []struct {
		name    string
		args    args
		want    *Set[E]
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				in:   []E{1, 2, 3},
				vals: []E{3},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					1: {},
					2: {},
				},
			},
			wantErr: false,
		},
		{
			name: "remove nonexistent",
			args: args{
				in:   []E{1, 2, 3},
				vals: []E{4},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					1: {},
					2: {},
					3: {},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.args.in...)
			s.Delete(tt.args.vals...)
			if !reflect.DeepEqual(s, tt.want) && !tt.wantErr {
				t.Errorf("Remove() = %v, want %v", s, tt.want)
			}
		})
	}
}

func Test_set_Difference(t *testing.T) {
	type E = int
	type args struct {
		a []E
		b []E
	}
	tests := []struct {
		name string
		args args
		want *Set[E]
	}{
		{
			name: "ok",
			args: args{
				a: []E{1, 2, 3, 4},
				b: []E{3, 4},
			},
			want: New([]E{1, 2}...),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sa := New(tt.args.a...)
			sb := New(tt.args.b...)
			d := sa.Difference(sb)
			for _, v := range tt.want.Members() {
				if !d.Contains(v) {
					t.Errorf("Difference() = %v does not contain %v", d, v)
				}
			}
			if tt.want.Size() != d.Size() {
				t.Errorf("Difference() = %v has unexpected members", d)
			}
		})
	}
}

func Test_set_Disjoint(t *testing.T) {
	type E = int
	type args struct {
		a, b []E
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "disjoint",
			args: args{
				a: []E{1, 2},
				b: []E{3, 4},
			},
			want: true,
		},
		{
			name: "not disjoint",
			args: args{
				a: []E{1, 2},
				b: []E{2, 3},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.args.a...)
			b := New(tt.args.b...)
			dis := a.Disjoint(b)
			if dis && !tt.want {
				t.Errorf("Disjoint() = %v, want %v", dis, tt.want)
			}
		})
	}
}

func Test_set_Empty(t *testing.T) {
	type E = int
	type args struct {
		in []E
	}
	tests := []struct {
		name      string
		args      args
		wantEmpty bool
	}{
		{
			name: "not empty",
			args: args{
				in: []E{1, 2, 3},
			},
			wantEmpty: false,
		},
		{
			name: "empty",
			args: args{
				in: []E{},
			},
			wantEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.args.in...)
			if s.Empty() && !tt.wantEmpty {
				t.Errorf("Empty() = %v, want %v", s, tt.wantEmpty)
			}
		})
	}
}

func Test_set_Intersection(t *testing.T) {
	type E = int
	type args struct {
		a, b []E
	}
	tests := []struct {
		name    string
		args    args
		want    *Set[E]
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				a: []E{1, 2},
				b: []E{2, 3},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					2: {},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				a: []E{1, 2},
				b: []E{2, 3},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					4: {},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.args.a...)
			b := New(tt.args.b...)
			inter := a.Intersection(b)
			if !reflect.DeepEqual(inter, tt.want) && !tt.wantErr {
				t.Errorf("Intersection() = %v, want %v", inter, tt.want)
			}
		})
	}
}

func Test_set_Members(t *testing.T) {
	type E = int
	type args struct {
		in []E
	}
	tests := []struct {
		name    string
		args    args
		want    []E
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				in: []E{1, 2},
			},
			want:    []E{1, 2},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				in: []E{1, 2},
			},
			want:    []E{3, 4},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.args.in...)
			members := a.Members()
			for _, want := range tt.want {
				found := false
				for _, got := range members {
					if got == want {
						found = true
						break
					}
				}
				if !found && !tt.wantErr {
					t.Errorf("Members(): wanted value %v not in members %v", want, members)
				}
			}
		})
	}
}

func Test_set_String(t *testing.T) {
	type E = int
	type args struct {
		in []E
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "ok",
			args: args{
				in: []E{1, 2},
			},
			want: []string{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.args.in...)
			s := a.String()
			for _, v := range tt.want {
				if !strings.Contains(s, v) {
					t.Errorf("String(): wanted %v got %v", v, s)
				}
			}
		})
	}
}

func Test_set_Union(t *testing.T) {
	type E = int
	type args struct {
		a, b []E
	}
	tests := []struct {
		name   string
		args   args
		want   *Set[E]
		wanted bool
	}{
		{
			name: "ok",
			args: args{
				a: []E{1, 2},
				b: []E{2, 3},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					1: {},
					2: {},
					3: {},
				},
			},
			wanted: true,
		},
		{
			name: "error",
			args: args{
				a: []E{1, 2},
				b: []E{2, 3},
			},
			want: &Set[E]{
				values: map[E]struct{}{
					4: {},
				},
			},
			wanted: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := New(tt.args.a...)
			b := New(tt.args.b...)
			union := a.Union(b)
			if reflect.DeepEqual(union, tt.want) != tt.wanted {
				t.Errorf("Union() = %v, want %v", union, tt.want)
			}
		})
	}
}
