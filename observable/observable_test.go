package observable

import (
	"testing"
)

func TestDetachObserver(t *testing.T) {
	type args struct {
		topic string
		name  string
		state map[string]Observable
	}
	var (
		f1 = func(func() interface{}) {}
		f2 = func(func() interface{}) {}
	)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				topic: "topic-one",
				name:  "obs-one",
				state: map[string]Observable{
					"topic-one": {
						getter: nil,
						observers: map[string]func(func() interface{}){
							"obs-one": f1,
							"obs-two": f2,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no-observable",
			args: args{
				topic: "topic-x",
				name:  "obs-three",
				state: map[string]Observable{
					"topic-one": {
						getter: nil,
						observers: map[string]func(func() interface{}){
							"obs-one": f1,
							"obs-two": f2,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no-observer",
			args: args{
				topic: "topic-one",
				name:  "obs-three",
				state: map[string]Observable{
					"topic-one": {
						getter: nil,
						observers: map[string]func(func() interface{}){
							"obs-one": f1,
							"obs-two": f2,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		state = tt.args.state
		t.Run(tt.name, func(t *testing.T) {
			if err := DetachObserver(tt.args.topic, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("DetachObserver() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteObservable(t *testing.T) {
	type args struct {
		topic string
		state map[string]Observable
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				topic: "topic-one",
				state: map[string]Observable{
					"topic-one": {
						getter:    nil,
						observers: map[string]func(func() interface{}){},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "no-observable",
			args: args{
				topic: "topic-ten",
				state: map[string]Observable{
					"topic-one": {
						getter:    nil,
						observers: map[string]func(func() interface{}){},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		state = tt.args.state
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteObservable(tt.args.topic); (err != nil) != tt.wantErr {
				t.Errorf("DeleteObservable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNotify(t *testing.T) {
	type args struct {
		topic string
		state map[string]Observable
	}
	var (
		f1 = func(func() interface{}) {}
	)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				topic: "topic-one",
				state: map[string]Observable{
					"topic-one": {
						getter: func() interface{} { return struct{}{} },
						observers: map[string]func(func() interface{}){
							"obs-one": f1,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				topic: "topic-two",
				state: make(map[string]Observable, 0),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		state = tt.args.state
		t.Run(tt.name, func(t *testing.T) {
			if err := Notify(tt.args.topic); (err != nil) != tt.wantErr {
				t.Errorf("Notify() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegisterObserver(t *testing.T) {
	type args struct {
		topic  string
		name   string
		update func(func() interface{})
		state  map[string]Observable
	}
	var (
		f1 = func(func() interface{}) {}
		f2 = func(func() interface{}) {}
	)
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantState map[string]Observable
	}{
		{
			name: "1",
			args: args{
				topic:  "topic-one",
				name:   "obs-one",
				update: f1,
				state:  make(map[string]Observable, 0),
			},
			wantErr: false,
			wantState: map[string]Observable{
				"topic-one": {
					getter: nil,
					observers: map[string]func(func() interface{}){
						"obs-one": f1,
					},
				},
			},
		},
		{
			name: "2",
			args: args{
				topic:  "topic-two",
				name:   "obs-two",
				update: f2,
				state: map[string]Observable{
					"topic-two": {
						getter: nil,
						observers: map[string]func(func() interface{}){
							"obs-one": f1,
						},
					},
				},
			},
			wantErr: false,
			wantState: map[string]Observable{
				"topic-two": {
					getter: nil,
					observers: map[string]func(func() interface{}){
						"obs-one": f1,
						"obs-two": f2,
					},
				},
			},
		},
		{
			name: "3",
			args: args{
				topic:  "topic-three",
				name:   "obs-three",
				update: f2,
				state: map[string]Observable{
					"topic-three": {
						getter:    nil,
						observers: nil,
					},
				},
			},
			wantErr: true,
			wantState: map[string]Observable{
				"topic-two": {
					getter:    nil,
					observers: nil,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state = tt.args.state
			err := RegisterObserver(tt.args.topic, tt.args.name, tt.args.update)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterObserver() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				_, exists := state[tt.args.topic]
				if !exists {
					t.Errorf("RegisterObserver() state[%s] does not exist", tt.args.topic)
				}
				_, exists = state[tt.args.topic].observers[tt.args.name]
				if !exists {
					t.Errorf("RegisterObserver() state[%s].observers[%s] does not exist", tt.args.topic, tt.args.name)
				}
				_, exists = state[tt.args.topic].observers[tt.args.name]
				if !exists {
					t.Errorf("RegisterObserver() state[%s].observers[%s] does not exist", tt.args.topic, tt.args.name)
				}
			}
		})
	}
}

func TestRegisterSubject(t *testing.T) {
	type args struct {
		topic string
		get   func() interface{}
		state map[string]Observable
	}
	var (
		f1 = func() interface{} { return nil }
		f2 = func(func() interface{}) {}
	)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				topic: "topic-one",
				get:   f1,
				state: make(map[string]Observable, 0),
			},
			wantErr: false,
		},
		{
			name: "existing-no-observer",
			args: args{
				topic: "topic-one",
				get:   f1,
				state: map[string]Observable{
					"topic-one": {
						getter: f1,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "existing-ok",
			args: args{
				topic: "topic-one",
				get:   f1,
				state: map[string]Observable{
					"topic-one": {
						getter: f1,
						observers: map[string]func(func() interface{}){
							"obs-one": f2,
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		state = tt.args.state
		t.Run(tt.name, func(t *testing.T) {
			err := RegisterSubject(tt.args.topic, tt.args.get)
			if (err != nil) != tt.wantErr {
				t.Errorf("RegisterSubject() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				s, exists := state[tt.args.topic]
				if !exists {
					t.Errorf("RegisterSubject() state[%s] does not exist", tt.args.topic)
				}
				if s.getter == nil {
					t.Errorf("RegisterSubject() getter for state[%s] is nil", tt.args.topic)
				}
			}
		})
	}
}
