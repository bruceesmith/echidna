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
	"context"
	"log/slog"
	"reflect"
	"testing"

	"github.com/bruceesmith/echidna/logger"
	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/urfave/cli/v3"
)

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
		ctx context.Context
		cmd *cli.Command
	}
	tests := []struct {
		name     string
		args     args
		wantCctx context.Context
		wantErr  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCctx, err := before(tt.args.ctx, tt.args.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("before() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotCctx, tt.wantCctx) {
				t.Errorf("before() = %v, want %v", gotCctx, tt.wantCctx)
			}
		})
	}
}

func Test_configure(t *testing.T) {
	type args struct {
		config        Configuration
		configLoaders []configLoader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := configure(tt.args.config, tt.args.configLoaders); (err != nil) != tt.wantErr {
				t.Errorf("configure() error = %v, wantErr %v", err, tt.wantErr)
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
			// wantValue: 77,
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
					Parser:   json.Parser(),
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
		// clArgs []string
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
			// clArgs = tt.args.line
			tt.args.command.Run(context.Background(), tt.args.line)
			log, _ := flag(tt.args.command, "log")
			t.Logf("json %v log %v trace %v", tt.args.command.Bool("json"), log, tt.args.command.StringSlice("trace"))
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
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := readConfig(tt.args.k, tt.args.sources...); (err != nil) != tt.wantErr {
				t.Errorf("readConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_printVersion(t *testing.T) {
	type args struct {
		cmd *cli.Command
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			printVersion(tt.args.cmd)
		})
	}
}

func TestRun(t *testing.T) {
	type args struct {
		ctx     context.Context
		command *cli.Command
		options []Option
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Run(tt.args.ctx, tt.args.command, tt.args.options...)
		})
	}
}

type config struct {
	I int
}

func (c *config) Validate() error { return nil }

func TestWithConfiguration(t *testing.T) {
	var cfg config

	type args struct {
		config Configuration
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				config: &cfg,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConfiguration(tt.args.config); got == nil {
				t.Errorf("WithConfiguration() = nil")
			}
		})
	}
}
