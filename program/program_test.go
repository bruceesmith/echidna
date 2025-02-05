// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

/*
Package program builds upon the Github packages knadh/koanf and urfave/cli/v3 to make it extremely simple to use the
features of those two excellent packages in concert.

Every program using program will expose a standard set of command-line flags (--json, --log, --trace, --verbose) in
addition to the standard flags provided by urfave/cli/v3 (--help and --version).

If a configuration struct is provided to the Run() function, then a further command=line flag (--config) is added to
provide the source(s) of values for fields in the struct.
*/
package program

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"testing"
	"unsafe"

	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/set"
	"github.com/knadh/koanf"
	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/urfave/cli/v3"
)

func Test_flagset_Delete(t *testing.T) {
	type fields struct {
		all   map[string]cli.Flag
		inuse *set.Set[string]
	}
	type args struct {
		name string
	}
	one := cli.BoolFlag{Name: "one"}
	two := cli.BoolFlag{Name: "two"}
	three := cli.IntFlag{Name: "three"}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				all: map[string]cli.Flag{
					"one":   &one,
					"two":   &two,
					"three": &three,
				},
				inuse: set.New("one", "two", "three"),
			},
			args: args{
				name: "three",
			},
		},
		{
			name: "non-existent",
			fields: fields{
				all: map[string]cli.Flag{
					"one":   &one,
					"two":   &two,
					"three": &three,
				},
				inuse: set.New("one", "two", "three"),
			},
			args: args{
				name: "four",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flagset{
				all:   tt.fields.all,
				inuse: tt.fields.inuse,
			}
			fs.Delete(tt.args.name)
			if found := tt.fields.inuse.Contains(tt.args.name); found {
				t.Errorf("flagset.Delete() failed to remove %s", tt.args.name)
			}
		})
	}
}

func Test_flagset_InUse(t *testing.T) {
	type fields struct {
		all   map[string]cli.Flag
		inuse *set.Set[string]
	}
	one := cli.BoolFlag{Name: "one"}
	two := cli.BoolFlag{Name: "two"}
	three := cli.IntFlag{Name: "three"}
	four := cli.StringFlag{Name: "four"}
	tests := []struct {
		name   string
		fields fields
		want   []cli.Flag
	}{
		{
			name: "ok",
			fields: fields{
				all: map[string]cli.Flag{
					"one":   &one,
					"two":   &two,
					"three": &three,
				},
				inuse: set.New("one", "three"),
			},
			want: []cli.Flag{
				&one,
				&three,
			},
		},
		{
			name: "non-existent",
			fields: fields{
				all: map[string]cli.Flag{
					"one":   &one,
					"two":   &two,
					"three": &three,
				},
				inuse: set.New("one", "three"),
			},
			want: []cli.Flag{
				&one,
				&four,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flagset{
				all:   tt.fields.all,
				inuse: tt.fields.inuse,
			}
			got := fs.InUse()
			if len(got) != len(tt.want) {
				t.Errorf("flagset.InUse() = wanted %d flag names, got %d", len(tt.want), len(got))
			}
			for _, f := range fs.inuse.Members() {
				expected := false
				name := ""
			loop:
				for _, used := range got {
					for _, name := range used.Names() {
						if f == name {
							expected = true
							break loop
						}
					}
				}
				if !expected {
					t.Errorf("flagset.InUse says %s is still in use", name)
				}
			}
		})
	}
}

func Test_flagset_Len(t *testing.T) {
	type fields struct {
		all   map[string]cli.Flag
		inuse *set.Set[string]
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "ok",
			fields: fields{
				all: map[string]cli.Flag{
					"one": &cli.BoolFlag{},
					"two": &cli.BoolFlag{},
				},
				inuse: set.New("one", "two"),
			},
			want: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := flagset{
				all:   tt.fields.all,
				inuse: tt.fields.inuse,
			}
			if got := fs.Len(); got != tt.want {
				t.Errorf("flagset.Len() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addCommand(t *testing.T) {
	type args struct {
		cmd     *cli.Command
		command *cli.Command
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "ok",
			args: args{
				cmd: &cli.Command{
					Commands: []*cli.Command{
						{
							Name: "one",
						},
					},
				},
				command: &cli.Command{
					Commands: []*cli.Command{
						{
							Name: "two",
						},
					},
				},
			},
			want: []string{
				"one",
				"two",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addCommand(tt.args.cmd, tt.args.command)
		})
	outer:
		for _, w := range tt.want {
			found := false
			for _, c := range tt.args.cmd.Commands {
				if c.Name == w {
					found = true
					break outer
				}
			}
			if !found {
				t.Errorf("addCommand() %s not found", w)
			}
		}
	}
}

func Test_addFlags(t *testing.T) {
	type args struct {
		cmd   *cli.Command
		flags []cli.Flag
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "ok",
			args: args{
				cmd: &cli.Command{
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "b",
						},
					},
				},
				flags: []cli.Flag{
					&cli.BoolFlag{
						Name: "bb",
					},
				},
			},
			want: []string{
				"b",
				"bb",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addFlags(tt.args.cmd, tt.args.flags)
		})
	outer:
		for _, w := range tt.want {
			for _, f := range tt.args.cmd.Flags {
				found := false
				for _, n := range f.Names() {
					if w == n {
						found = true
						break outer
					}
				}
				if !found {
					t.Errorf("addflags() %s not found", w)
				}
			}
		}
	}
}

func Test_before(t *testing.T) {
	type args struct {
		cmd  *cli.Command
		line []string
		cfg  Configurator
	}
	var (
		cfg   config
		fail1 vfail
	)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				cmd: &cli.Command{
					Name: "test",
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Before: before,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "b",
						},
						&cli.StringSliceFlag{
							Name: "config",
						},
					},
				},
				line: []string{"test", "--b", "--config", "testdata/test.yml"},
				cfg:  &cfg,
			},
			wantErr: false,
		},
		{
			name: "file-not-exit",
			args: args{
				cmd: &cli.Command{
					Name: "test",
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Before: before,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "b",
						},
						&cli.StringSliceFlag{
							Name: "config",
						},
					},
				},
				line: []string{"test", "--b", "--config", "testdata/does-not-exist.yml"},
				cfg:  &cfg,
			},
			wantErr: true,
		},
		{
			name: "logging-error",
			args: args{
				cmd: &cli.Command{
					Name: "test",
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Before: before,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "b",
						},
						&cli.StringSliceFlag{
							Name: "config",
						},
						&logger.LogLevelFlag{
							Name: "log",
						},
					},
				},
				line: []string{"test", "--b", "--config", "testdata/test.yml", "--log", "fred"},
				cfg:  &cfg,
			},
			wantErr: true,
		},
		{
			name: "validation-failed",
			args: args{
				cmd: &cli.Command{
					Name: "test",
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Before: before,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "b",
						},
						&cli.StringSliceFlag{
							Name: "config",
						},
					},
				},
				line: []string{"test", "--b", "--config", "testdata/test.yml"},
				cfg:  &fail1,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configuration = &cfg
			buf := &bytes.Buffer{}
			tt.args.cmd.Writer = buf
			tt.args.cmd.ErrWriter = buf
			err := tt.args.cmd.Run(context.Background(), tt.args.line)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("error %v", err)
				}
			}
			configuration = nil
		})
	}
}

func Test_configure(t *testing.T) {
	type args struct {
		config        Configurator
		configLoaders []configLoader
	}
	var (
		cfg config
		s1  = struct {
			I int `koanf:"i"`
		}{
			I: 33,
		}
		s2 = struct {
			I int `koanf:"i"`
		}{
			I: 22,
		}
	)
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				config: &cfg,
				configLoaders: []configLoader{
					{
						Provider: structs.Provider(s1, "koanf"),
						Parser:   nil,
						Options:  nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "read-error",
			args: args{
				config: &cfg,
				configLoaders: []configLoader{
					{
						Provider: structs.Provider(s1, "koanf"),
						Parser:   kjson.Parser(),
						Options:  nil,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "validate-error",
			args: args{
				config: &cfg,
				configLoaders: []configLoader{
					{
						Provider: structs.Provider(s2, "koanf"),
						Parser:   nil,
						Options:  nil,
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := configure(tt.args.config, tt.args.configLoaders); err != nil {
				if !tt.wantErr {
					t.Errorf("configure() error = %v, wantErr %v", err, tt.wantErr)
				}
			}
		})
	}
}

func TestConfiguration(t *testing.T) {
	var (
		cfg = config{I: 33}
		s1  simple1
		s2  simple2
	)

	type args struct {
		config Configurator
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				config: &cfg,
			},
			wantErr: false,
		},
		{
			name: "not-a-pointer",
			args: args{
				config: s1,
			},
			wantErr: true,
		},
		{
			name: "not-a-struct",
			args: args{
				config: &s2,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Option
			if got = Configuration(tt.args.config); got == nil {
				t.Errorf("Configuration() = nil")
			}
			err := got()
			if err != nil && !tt.wantErr {
				t.Errorf("Configuration() err %v wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr {
				if configuration != tt.args.config {
					t.Errorf("Configuration() configuration %p not expected, want %p", configuration, tt.args.config)
				}
			}
			configuration = nil
		})
	}
}

func Test_expand(t *testing.T) {
	type args struct {
		slice []int
		size  int
	}
	ar := [4]int{1, 2, 3, 4}
	sl := ar[0:2]
	tests := []struct {
		name                 string
		args                 args
		wantCap              int
		wantLen              int
		wantDifferentAddress bool
	}{
		{
			name: "will-expand",
			args: args{
				slice: []int{9, 8},
				size:  2,
			},
			wantCap:              4,
			wantLen:              2,
			wantDifferentAddress: true,
		},
		{
			name: "will-not-expand",
			args: args{
				slice: sl,
				size:  1,
			},
			wantCap:              4,
			wantLen:              2,
			wantDifferentAddress: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV := expand(tt.args.slice, tt.args.size)
			if len(gotV) != tt.wantLen || cap(gotV) != tt.wantCap {
				t.Errorf("expand() = %v, want len/cap %v/%v got len/cap %v/%v", gotV, tt.wantLen, tt.wantCap, len(gotV), cap(gotV))
			}
			if (unsafe.SliceData(tt.args.slice) != unsafe.SliceData(gotV)) != tt.wantDifferentAddress {
				t.Errorf("expand() original %v expanded %v", unsafe.SliceData(tt.args.slice), unsafe.SliceData(gotV))
			}
		})
	}
}

func Test_flag(t *testing.T) {
	type args struct {
		cmd  *cli.Command
		name string
		line []string
	}
	tests := []struct {
		name      string
		args      args
		wantValue any
		wantFound bool
	}{
		{
			name: "set-and-found",
			args: args{
				cmd: &cli.Command{
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Name: "test",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "someflag",
							Value: "fred",
						},
					},
				},
				name: "someflag",
				line: []string{"test", "-someflag", "bill"},
			},
			wantValue: "bill",
			wantFound: true,
		},
		{
			name: "found-but-not-set",
			args: args{
				cmd: &cli.Command{
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Name: "test",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "someflag",
							Value: "fred",
						},
					},
				},
				name: "someflag",
				line: []string{"test"},
			},
			wantValue: "fred",
			wantFound: true,
		},
		{
			name: "custom-set-and-found",
			args: args{
				cmd: &cli.Command{
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Name: "test",
					Flags: []cli.Flag{
						&logger.LogLevelFlag{
							Name:  "log",
							Value: logger.LogLevel(slog.LevelWarn),
						},
					},
				},
				name: "log",
				line: []string{"test", "-log", "error"},
			},
			wantValue: logger.LogLevel(slog.LevelError),
			wantFound: true,
		},
		{
			name: "not-found",
			args: args{
				cmd: &cli.Command{
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Name: "test",
					Flags: []cli.Flag{
						&cli.StringFlag{
							Name:  "someotherflag",
							Value: "fred",
						},
					},
				},
				name: "someflag",
				line: []string{"test", "-someotherflag", "harry"},
			},
			wantValue: nil,
			wantFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.cmd.Action != nil {
				tt.args.cmd.Run(context.Background(), tt.args.line)
			}
			gotValue, gotFound := flag(tt.args.cmd, tt.args.name)
			if gotValue != tt.wantValue {
				t.Errorf("flag() gotValue = %v, want %v", gotValue, tt.wantValue)
			}
			if gotFound != tt.wantFound {
				t.Errorf("flag() gotFound = %v, want %v", gotFound, tt.wantFound)
			}
		})
	}
}

func Test_loaders(t *testing.T) {
	type args struct {
		paths []string
	}
	tests := []struct {
		name    string
		args    args
		want    []configLoader
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				paths: []string{
					"one.json",
					"two.yml",
				},
			},
			want: []configLoader{
				{
					Provider: file.Provider("one.json"),
					Parser:   kjson.Parser(),
					Options:  []koanf.Option{},
				},
				{
					Provider: file.Provider("two.yml"),
					Parser:   yaml.Parser(),
					Options:  []koanf.Option{},
				},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				paths: []string{
					"one.unknown",
				},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loaders(tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("loaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != len(tt.want) {
				t.Errorf("loaders() = got %v, want %v", len(got), len(tt.want))
			}
			for i, loader := range got {
				if loader.Provider == nil {
					t.Errorf("loaders() = got nil Provider at index %v", i)
				}
				if loader.Parser != tt.want[i].Parser {
					t.Errorf("loaders() = got Parser %v, want Parser %v at index %v", loader.Parser, tt.want[i].Parser, i)
				}
			}
		})
	}
}

func Test_logging(t *testing.T) {
	type args struct {
		command *cli.Command
		line    []string
	}
	var (
		tester = func(_ context.Context, _ *cli.Command) error {
			return nil
		}
	)
	tests := []struct {
		name      string
		args      args
		wantErr   bool
		wantLevel string
	}{
		{
			name: "ok",
			args: args{
				command: &cli.Command{
					Action: tester,
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "json",
						},
						&logger.LogLevelFlag{
							Name: "log",
						},
						&cli.StringSliceFlag{
							Name: "trace",
						},
					},
					Name: "test",
				},
				line: []string{"test", "-json", "-log", "warn", "-trace", "area1"},
			},
			wantErr:   false,
			wantLevel: "WARN",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.command.Run(context.Background(), tt.args.line)
			if err := logging(tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("logging() error = %v, wantErr %v", err, tt.wantErr)
			}
			if logger.Level() != tt.wantLevel {
				t.Errorf("logging() got level %v want level %v", logger.Level(), tt.wantLevel)
			}
		})
	}
}

func Test_readConfig(t *testing.T) {
	type args struct {
		k       *koanf.Koanf
		sources []configLoader
	}
	var s = struct {
		I int `koanf:"i"`
	}{
		I: 22,
	}
	tests := []struct {
		name       string
		args       args
		wantErr    bool
		wantStruct map[string]interface{}
	}{
		{
			name: "ok",
			args: args{
				k: koanf.New("."),
				sources: []configLoader{
					{
						Provider: structs.Provider(s, "koanf"),
						Parser:   nil,
						Options:  nil,
					},
				},
			},
			wantErr: false,
			wantStruct: map[string]interface{}{
				"i": 22,
			},
		},
		{
			name: "single-error",
			args: args{
				k: koanf.New("."),
				sources: []configLoader{
					{
						Provider: structs.Provider(s, "koanf"),
						Parser:   kjson.Parser(),
						Options:  nil,
					},
				},
			},
			wantErr: true,
			wantStruct: map[string]interface{}{
				"i": 22,
			},
		},
		{
			name: "multiple-error",
			args: args{
				k: koanf.New("."),
				sources: []configLoader{
					{
						Provider: structs.Provider(s, "koanf"),
						Parser:   kjson.Parser(),
						Options:  nil,
					},
					{
						Provider: structs.Provider(s, "koanf"),
						Parser:   kjson.Parser(),
						Options:  nil,
					},
				},
			},
			wantErr: true,
			wantStruct: map[string]interface{}{
				"i": 22,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := readConfig(tt.args.k, tt.args.sources...)
			if err != nil {
				if !tt.wantErr {
					t.Errorf("readConfig() error = %v, wantErr %v", err, tt.wantErr)
				}
			} else {
				if !reflect.DeepEqual(tt.args.k.Raw(), tt.wantStruct) {
					t.Errorf("readConfig() want %v got %v", tt.wantStruct, tt.args.k.Raw())
				}
			}
		})
	}
}

func Test_printVersion(t *testing.T) {
	type args struct {
		cmd *cli.Command
	}
	tests := []struct {
		name    string
		args    args
		jason   bool
		verbose bool
	}{
		{
			name: "simple",
			args: args{
				cmd: &cli.Command{
					Name:        "test",
					Description: "a test command",
					Action: func(cts context.Context, cmd *cli.Command) error {
						return nil
					},
					Version: "1",
				},
			},
			jason:   false,
			verbose: false,
		},
		{
			name: "json",
			args: args{
				cmd: &cli.Command{
					Name:        "test",
					Description: "a test command",
					Action: func(cts context.Context, cmd *cli.Command) error {
						return nil
					},
					Version: "1",
				},
			},
			jason:   true,
			verbose: false,
		},
		{
			name: "verbose",
			args: args{
				cmd: &cli.Command{
					Name:        "test",
					Description: "a test command",
					Action: func(cts context.Context, cmd *cli.Command) error {
						return nil
					},
					Version: "1",
				},
			},
			jason:   false,
			verbose: true,
		},
		{
			name: "json-and-verbose",
			args: args{
				cmd: &cli.Command{
					Name:        "test",
					Description: "a test command",
					Action: func(cts context.Context, cmd *cli.Command) error {
						return nil
					},
					Version: "1",
				},
			},
			jason:   true,
			verbose: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []string{"test", "--version"}
			if tt.jason {
				args = append(args, "--json")
			}
			if tt.verbose {
				args = append(args, "--verbose")
			}

			flags.Delete("config")
			buf := &bytes.Buffer{}
			tt.args.cmd.Writer = buf
			addFlags(tt.args.cmd, flags.InUse())
			cli.VersionPrinter = printVersion
			tt.args.cmd.Run(context.Background(), args)
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		command *cli.Command
		options []Option
		line    []string
	}
	var cfg config
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				command: &cli.Command{
					Name: "test",
					Action: func(context.Context, *cli.Command) error {
						return nil
					},
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "b",
						},
					},
				},
				options: []Option{
					Configuration(&cfg),
				},
				line: []string{"test", "--b"},
			},
			wantErr: false,
		},
		{
			name: "error",
			args: args{
				command: &cli.Command{
					Name: "test",
					Action: func(context.Context, *cli.Command) error {
						return errors.New("Run() error")
					},
					Flags: []cli.Flag{
						&cli.BoolFlag{
							Name: "b",
						},
					},
				},
				options: []Option{
					Configuration(&cfg),
				},
				line: []string{"test", "--b"},
			},
			wantErr: true,
		},
	}
	before := flags.inuse.Members()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := &bytes.Buffer{}
			tt.args.command.Writer = buf
			noOsExit = true
			os.Args = tt.args.line
			Run(context.Background(), tt.args.command, tt.args.options...)
			noOsExit = false
		})
	}
	flags.inuse = set.New(before...)
}

type config struct {
	I int `koanf:"i"`
}

func (c *config) Validate() error {
	if c.I != 33 {
		return errors.New("I must be 33")
	}
	return nil
}

type simple1 int

func (i simple1) Validate() error {
	return nil
}

type simple2 int

type vfail int

func (v vfail) Validate() error {
	return fmt.Errorf("validation failed")
}

func (i *simple2) Validate() error {
	return nil
}
func TestNoDefaultFlags(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Option
			if got = NoDefaultFlags(); got == nil {
				t.Errorf("NoDefaultFlags() returned nil ")
			}
			got()
			if flags.inuse.Contains("json") || flags.inuse.Contains("log") || flags.inuse.Contains("trace") || flags.inuse.Contains("verbose") {
				t.Errorf("NoDefaultFlags() unexpected in-use flags = %v", flags.inuse.Members())
			}

		})
	}
}

func TestNoJSON(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Option
			if got = NoJSON(); got == nil {
				t.Errorf("NoJson() returned nil ")
			}
			got()
			if flags.inuse.Contains("json") {
				t.Error("NoJson failed to remove the json flag")
			}
		})
	}
}

func TestNoLog(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Option
			if got = NoLog(); got == nil {
				t.Errorf("NoLog() returned nil ")
			}
			got()
			if flags.inuse.Contains("log") {
				t.Error("NoLog failed to remove the log flag")
			}
		})
	}
}

func TestNoTrace(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Option
			if got = NoTrace(); got == nil {
				t.Errorf("NoTrace() returned nil ")
			}
			got()
			if flags.inuse.Contains("trace") {
				t.Error("NoTrace failed to remove the trace flag")
			}
		})
	}
}

func TestNoVerbose(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "ok",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Option
			if got = NoVerbose(); got == nil {
				t.Errorf("NoVerbose() returned nil ")
			}
			got()
			if flags.inuse.Contains("verbose") {
				t.Error("NoVerbose failed to remove the verbose flag")
			}
		})
	}
}
