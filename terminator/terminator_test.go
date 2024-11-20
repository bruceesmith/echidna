// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package terminator

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	type args struct {
		delta int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				delta: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefault(New())
			go func() {
				Add(tt.args.delta)
				for range tt.args.delta {
					Done()
				}
			}()
			Wait()
		})
	}
}

func TestDefault(t *testing.T) {
	tests := []struct {
		name string
		want *Terminator
	}{
		{
			name: "ok",
			want: defaultTerminator.Load(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Default(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Default() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDone(t *testing.T) {
	tests := []struct {
		name  string
		delta int
	}{
		{
			name:  "ok",
			delta: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefault(New())
			go func() {
				Add(tt.delta)
				for range tt.delta {
					Done()
				}
			}()
			Wait()
		})
	}
}

func TestShutDown(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefault(New())
			Stop()
			select {
			case <-ShutDown():
			default:
				t.Error("ShutDown() did not get signal")
			}
		})
	}
}

func TestSetDefault(t *testing.T) {
	type args struct {
		l *Terminator
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefault(tt.args.l)
		})
	}
}

func TestShuttingDown(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefault(New())
			Stop()
			if !ShuttingDown() {
				t.Error("ShuttingDown() should return true")
			}
		})
	}
}

func TestStop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefault(New())
			Stop()
			select {
			case <-ShutDown():
				if !Default().shuttingDown {
					t.Error("Stop() did not set 'shuttingDown'")
				}
			default:
				t.Error("Stop() did not get signal")
			}
		})
	}
}

func TestWait(t *testing.T) {
	tests := []struct {
		name  string
		delta int
	}{
		{
			name:  "ok",
			delta: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetDefault(New())
			go func() {
				Add(tt.delta)
				for range tt.delta {
					Done()
				}
			}()
			Wait()
		})
	}
}

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Terminator
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTerminator_Add(t *testing.T) {
	type args struct {
		delta int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				delta: 2,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := New()
			go func() {
				tm.Add(tt.args.delta)
				for range tt.args.delta {
					tm.Done()
				}
			}()
			tm.Wait()
		})
	}
}

func TestTerminator_Done(t *testing.T) {
	tests := []struct {
		name  string
		delta int
	}{
		{
			name:  "ok",
			delta: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := New()
			go func() {
				tm.Add(tt.delta)
				for range tt.delta {
					tm.Done()
				}
			}()
			tm.Wait()
		})
	}
}

func TestTerminator_ShutDown(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := New()
			tm.Stop()
			select {
			case <-tm.ShutDown():
			default:
				t.Error("ShutDown() did not get signal")
			}
		})
	}
}

func TestTerminator_ShuttingDown(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := New()
			tm.Stop()
			if !tm.ShuttingDown() {
				t.Error("ShuttingDown() should return true")
			}
		})
	}
}

func TestTerminator_Stop(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := New()
			tm.Stop()
			select {
			case <-tm.ShutDown():
				if !tm.shuttingDown {
					t.Error("Stop() did not set 'shuttingDown'")
				}
			default:
				t.Error("Stop() did not get signal")
			}
		})
	}
}

func TestTerminator_Wait(t *testing.T) {
	tests := []struct {
		name  string
		delta int
	}{
		{
			name:  "ok",
			delta: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tm := New()
			go func() {
				tm.Add(tt.delta)
				for range tt.delta {
					tm.Done()
				}
			}()
			tm.Wait()
		})
	}
}
