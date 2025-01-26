// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package stack

import (
	"sync"
	"testing"
)

func TestStack_IsEmpty(t *testing.T) {
	type fields struct {
		top  *element[int]
		mut  *sync.Mutex
		size int
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "empty",
			fields: fields{
				mut: &sync.Mutex{},
			},
			want: true,
		},
		{
			name: "not-empty",
			fields: fields{
				top:  &element[int]{},
				mut:  &sync.Mutex{},
				size: 1,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack[int]{
				top:  tt.fields.top,
				mut:  tt.fields.mut,
				size: tt.fields.size,
			}
			if got := s.IsEmpty(); got != tt.want {
				t.Errorf("Stack.IsEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStack_Peek(t *testing.T) {
	type fields struct {
		top  *element[int]
		mut  *sync.Mutex
		size int
	}
	tests := []struct {
		name      string
		fields    fields
		wantValue int
		wantOK    bool
	}{
		{
			name: "empty",
			fields: fields{
				mut: &sync.Mutex{},
			},
			wantOK: false,
		},
		{
			name: "not-empty",
			fields: fields{
				top: &element[int]{
					value: 11,
					next:  nil,
				},
				mut:  &sync.Mutex{},
				size: 1,
			},
			wantValue: 11,
			wantOK:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack[int]{
				top:  tt.fields.top,
				mut:  tt.fields.mut,
				size: tt.fields.size,
			}
			top := s.top
			siz := s.size
			gotValue, gotOK := s.Peek()
			if gotValue != tt.wantValue {
				t.Errorf("Stack.Peek() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOK != tt.wantOK {
				t.Errorf("Stack.Peek() gotEmpty = %v, want %v", gotOK, tt.wantOK)
			}
			if top != s.top {
				t.Errorf("Stack.Peek() got top = %v, want %v", s.top, top)
			}
			if siz != s.size {
				t.Errorf("Stack.Peek() got size = %v, want %v", s.size, siz)
			}
		})
	}
}

func TestStack_Pop(t *testing.T) {
	type fields struct {
		top  *element[int]
		mut  *sync.Mutex
		size int
	}
	tests := []struct {
		name      string
		fields    fields
		wantValue int
		wantOK    bool
	}{
		{
			name: "empty",
			fields: fields{
				mut: &sync.Mutex{},
			},
			wantValue: 0,
			wantOK:    false,
		},
		{
			name: "not-empty",
			fields: fields{
				top: &element[int]{
					value: 23,
					next:  nil,
				},
				mut:  &sync.Mutex{},
				size: 1,
			},
			wantValue: 23,
			wantOK:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack[int]{
				top:  tt.fields.top,
				mut:  tt.fields.mut,
				size: tt.fields.size,
			}
			size := s.size
			var gotValue int
			var gotOK bool
			gotValue, gotOK = s.Pop()
			if gotValue != tt.wantValue {
				t.Errorf("Stack.Pop() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotOK != tt.wantOK {
				t.Errorf("Stack.Pop() gotEmpty = %v, want %v", gotOK, tt.wantOK)
			}
			if gotOK && s.size != size-1 {
				t.Errorf("Stack.Pop() got size = %v, want %v", s.size, size-1)
			}
		})
	}
}

func TestStack_Push(t *testing.T) {
	type fields struct {
		top  *element[int]
		mut  *sync.Mutex
		size int
	}
	type args struct {
		v int
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		wantSize     int
		wantTopValue int
	}{
		{
			name: "start-empty",
			fields: fields{
				mut: &sync.Mutex{},
			},
			args: args{
				v: 77,
			},
			wantSize:     1,
			wantTopValue: 77,
		},
		{
			name: "start-not-empty",
			fields: fields{
				top: &element[int]{
					value: 44,
					next:  nil,
				},
				mut:  &sync.Mutex{},
				size: 1,
			},
			args: args{
				v: 66,
			},
			wantSize:     2,
			wantTopValue: 66,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack[int]{
				top:  tt.fields.top,
				mut:  tt.fields.mut,
				size: tt.fields.size,
			}
			s.Push(tt.args.v)
			if tt.wantSize != s.size {
				t.Errorf("Stack.Push() got size = %v, want %v", s.size, tt.wantSize)
			}
			if s.top.value != tt.wantTopValue {
				t.Errorf("Stack.Push() got top value = %v, want %v", s.top.value, tt.wantTopValue)
			}
		})
	}
}

func TestStack_Size(t *testing.T) {
	type fields struct {
		top  *element[int]
		mut  *sync.Mutex
		size int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "empty",
			fields: fields{
				mut: &sync.Mutex{},
			},
			want: 0,
		},
		{
			name: "not-empty",
			fields: fields{
				top: &element[int]{
					value: 99,
					next:  nil,
				},
				mut:  &sync.Mutex{},
				size: 1,
			},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Stack[int]{
				top:  tt.fields.top,
				mut:  tt.fields.mut,
				size: tt.fields.size,
			}
			if got := s.Size(); got != tt.want {
				t.Errorf("Stack.Size() = %v, want %v", got, tt.want)
			}
		})
	}
}
