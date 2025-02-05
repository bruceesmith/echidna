// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package program

import (
	"reflect"
	"testing"

	"github.com/urfave/cli/v3"
	"github.com/urfave/sflags"
)

func TestDescTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				tag: "flagtag",
			},
			want: "flagtag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DescTag(tt.args.tag)
			got()
			if opts.descTag != tt.want {
				t.Errorf("DescTag() = %v, want %v", opts.descTag, tt.want)
			}
		})
	}
}

func TestEnvDivider(t *testing.T) {
	type args struct {
		divider string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				divider: "x",
			},
			want: "x",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnvDivider(tt.args.divider)
			got()
			if opts.envDivider != tt.want {
				t.Errorf("EnvDivider() = %v, want %v", opts.envDivider, tt.want)
			}
		})
	}
}

func TestEnvPrefix(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				prefix: "envo",
			},
			want: "envo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EnvPrefix(tt.args.prefix)
			got()
			if opts.envPrefix != tt.want {
				t.Errorf("EnvPrefix() = %v, want %v", opts.envPrefix, tt.want)
			}
		})
	}
}

func TestFlagDivider(t *testing.T) {
	type args struct {
		divider string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				divider: ".",
			},
			want: ".",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlagDivider(tt.args.divider)
			got()
			if opts.divider != tt.want {
				t.Errorf("FlagDivider() = %v, want %v", opts.divider, tt.want)
			}
		})
	}
}

func TestFlagTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				tag: "flag",
			},
			want: "flag",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FlagTag(tt.args.tag)
			got()
			if opts.tag != tt.want {
				t.Errorf("FlagTag() = %v, want %v", opts.tag, tt.want)
			}
		})
	}
}

func TestFlatten(t *testing.T) {
	type args struct {
		flatten bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ok",
			args: args{
				flatten: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Flatten(tt.args.flatten)
			got()
			if opts.flatten != tt.want {
				t.Errorf("Flatten() = %v, want %v", opts.flatten, tt.want)
			}
		})
	}
}

func TestPrefix(t *testing.T) {
	type args struct {
		prefix string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ok",
			args: args{
				prefix: "config",
			},
			want: "config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Prefix(tt.args.prefix)
			got()
			if opts.prefix != tt.want {
				t.Errorf("Prefix() = %v, want %v", opts.prefix, tt.want)
			}
		})
	}
}

func TestValidator(t *testing.T) {
	type args struct {
		val sflags.ValidateFunc
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				val: func(_ string, f reflect.StructField, _ any) error {
					return nil
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Validator(tt.args.val)
			got()
			if opts.validator == nil {
				t.Error("Validator() = opts.validator is not set")
			}
		})
	}
}

func Test_applyFlagOverrides(t *testing.T) {
	type args struct {
		names       []string
		flagI, cfgI int
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				names: []string{"i"},
				flagI: 77,
				cfgI:  33,
			},
		},
		{
			name: "no-update",
			args: args{
				names: []string{"i"},
				flagI: 22,
				cfgI:  81,
			},
		},
	}
	opts = options{
		descTag:    "desc",
		divider:    "-",
		envDivider: "_",
		envPrefix:  "",
		tag:        "flag",
		flatten:    false,
		prefix:     "",
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := config{I: tt.args.flagI}
			bdr, _ := newFlagBinder(&cfg)
			cfg.I = tt.args.cfgI
			applyFlagOverrides(tt.args.names, bdr)
			if !reflect.DeepEqual(bdr.clone, &cfg) {
				t.Errorf("applyFlagOverrides() updates not applied %+v %+v", bdr.clone, &cfg)
			}
		})
	}
}

func Test_clone(t *testing.T) {
	type args struct {
		in Configurator
	}
	var cfg config
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				in: &cfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := clone(tt.args.in)
			if reflect.ValueOf(got) == reflect.ValueOf(tt.args.in) {
				t.Errorf("clone() addresses are the same")
			}
			if got == tt.args.in {
				t.Errorf("clone() addresses are the same")
			}
		})
	}
}

func Test_fieldMap(t *testing.T) {
	type args struct {
		v Configurator
	}
	var (
		cfg config
		i   simple1
	)
	tests := []struct {
		name    string
		args    args
		want    map[string]reflect.Value
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				v: &cfg,
			},
			want: map[string]reflect.Value{
				"i": reflect.ValueOf(cfg.I),
			},
			wantErr: false,
		},
		{
			name: "nil pointer",
			args: args{
				v: nil,
			},
			want: map[string]reflect.Value{
				"i": reflect.ValueOf(cfg.I),
			},
			wantErr: true,
		},
		{
			name: "pointer_to_int",
			args: args{
				v: &i,
			},
			want: map[string]reflect.Value{
				"i": reflect.ValueOf(cfg.I),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts = options{
				divider: "-",
				tag:     "flag",
				flatten: false,
				prefix:  "",
			}
			got, err := fieldMap(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("fieldMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			for k := range tt.want {
				found := false
				for kk := range got {
					if k == kk {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("fieldMap() expected key = %v, got %v", k, got)
				}
			}
		})
	}
}

type config1 struct {
	C1 config
	s1 simple1
}

func Test_mapFields(t *testing.T) {
	type args struct {
		v      reflect.Value
		prefix string
		result *map[string]reflect.Value
	}
	var (
		cfg      config
		cfg1     config1
		fieldmap = make(map[string]reflect.Value)
	)
	tests := []struct {
		name string
		args args
	}{
		{
			name: "simple_config",
			args: args{
				v:      reflect.ValueOf(&cfg).Elem(),
				prefix: "",
				result: &fieldmap,
			},
		},
		{
			name: "nested_config",
			args: args{
				v:      reflect.ValueOf(&cfg1).Elem(),
				prefix: "",
				result: &fieldmap,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mapFields(tt.args.v, tt.args.prefix, tt.args.result)
		})
	}
}

func Test_newFlagBinder(t *testing.T) {
	var cfg config
	tests := []struct {
		name    string
		config  Configurator
		wantB   binder
		wantErr bool
	}{
		{
			name:    "ok",
			config:  &cfg,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotB, err := newFlagBinder(tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("newFlagBinder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(gotB.configFields) != len(gotB.cloneFields) {
				t.Errorf("newFlagBinder() = %v, want %v", len(gotB.cloneFields), len(gotB.configFields))
			}
		})
	}
}

func TestConfigFlags(t *testing.T) {
	type args struct {
		configs []Configurator
		ops     []FlagOption
		param   any
	}
	var (
		cfg config
		cmd = cli.Command{}
	)
	tests := []struct {
		name      string
		args      args
		wantFlags []string
		wantErr   bool
	}{
		{
			name: "ok",
			args: args{
				configs: []Configurator{&cfg},
				param:   &cmd,
			},
			wantFlags: []string{"i"},
			wantErr:   false,
		},
		{
			name: "expected-not-found",
			args: args{
				configs: []Configurator{&cfg},
				param:   &cmd,
			},
			wantFlags: []string{"ii"},
			wantErr:   true,
		},
		{
			name: "no-configs",
			args: args{
				configs: []Configurator{},
				param:   &cmd,
			},
			wantErr: true,
		},
		{
			name: "not-a-command",
			args: args{
				configs: []Configurator{&cfg},
				param:   cfg,
			},
			wantErr: true,
		},
		{
			name: "different-tag",
			args: args{
				configs: []Configurator{&cfg},
				ops:     []FlagOption{FlagTag("fred")},
				param:   &cmd,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConfigFlags(tt.args.configs, &cmd, tt.args.ops...)
			err := got(tt.args.param)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("ConfigFlags() = %v, want %v %v", err, tt.wantErr, reflect.TypeOf(cmd))
				}
			}
			for _, wfn := range tt.wantFlags {
				found := false
			loop:
				for _, gf := range cmd.Flags {
					for _, nam := range gf.Names() {
						if wfn == nam {
							found = true
							break loop
						}
					}
				}
				if !found && !tt.wantErr {
					t.Errorf("ConfigFlags() = expected flag %v not found", wfn)
				}
			}
		})
	}
}
