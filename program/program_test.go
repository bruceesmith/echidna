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
	"reflect"
	"testing"

	"github.com/knadh/koanf"
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

func Test_logging(t *testing.T) {
	type args struct {
		command *cli.Command
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
			if err := logging(tt.args.command); (err != nil) != tt.wantErr {
				t.Errorf("logging() error = %v, wantErr %v", err, tt.wantErr)
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
	}
	tests := []struct {
		name      string
		args      args
		wantValue string
		wantFound bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := loaders(tt.args.paths)
			if (err != nil) != tt.wantErr {
				t.Errorf("loaders() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("loaders() = %v, want %v", got, tt.want)
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

func TestWithConfiguration(t *testing.T) {
	type args struct {
		config Configuration
	}
	tests := []struct {
		name string
		args args
		want Option
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WithConfiguration(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("WithConfiguration() = %v, want %v", got, tt.want)
			}
		})
	}
}
