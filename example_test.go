// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package echidna

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/urfave/cli/v3"
)

type configExample struct {
	I int
}

func (c *configExample) Validate() error { return nil }

func ExampleRun_basic() {
	// The most basic example
	//
	var cmd = &cli.Command{
		Name: "runbasic",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
	os.Args = []string{"runbasic", "h"}
	Run(context.Background(), cmd, NoDefaultFlags())
	// Output:
	// NAME:
	//    runbasic - A new cli application
	//
	// USAGE:
	//    runbasic [global options] [command [command options]]
	//
	// COMMANDS:
	//    version, v  print the version
	//    help, h     Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --help, -h  show help
}

func ExampleRun_version() {
	// Include a Version field in the Command
	//
	var cmd = &cli.Command{
		Name: "runversion",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
		Version: "1",
	}
	os.Args = []string{"runversion", "h"}
	Run(context.Background(), cmd, NoDefaultFlags())
	// Output:
	// NAME:
	//    runversion - A new cli application
	//
	// USAGE:
	//    runversion [global options] [command [command options]]
	//
	// VERSION:
	//    1
	//
	// COMMANDS:
	//    version, v  print the version
	//    help, h     Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --help, -h     show help
	//    --version, -v  print the version
}

func ExampleRun_action() {
	// Include an Action function
	//
	var cmd = &cli.Command{
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("hello")
			return nil
		},
		Name:    "runaction",
		Version: "1",
	}
	os.Args = []string{"runaction"}
	Run(context.Background(), cmd)
	// Output:
	// hello
}

func ExampleRun_customflag() {
	// Include a custom flag
	//
	var cmd = &cli.Command{
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("hello", cmd.Int("i"))
			return nil
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "i",
				Usage: "An integer",
			},
		},
		Name:    "runcustomflag",
		Version: "1",
	}
	os.Args = []string{"runcustomflag", "-i", "22"}
	Run(context.Background(), cmd, NoDefaultFlags())
	// Output:
	// hello 22
}

func ExampleRun_flagwithdefault() {
	// Include a custom flag with a default value
	//
	var cmd = &cli.Command{
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("hello")
			return nil
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "i",
				Usage: "An integer",
				Value: 22,
			},
		},
		Name:    "runflagwithdefault",
		Version: "1",
	}
	os.Args = []string{"runflagwithdefault", "--help"}
	Run(context.Background(), cmd, NoDefaultFlags())
	// Output:
	// NAME:
	//    runflagwithdefault - A new cli application
	//
	// USAGE:
	//    runflagwithdefault [global options] [command [command options]]
	//
	// VERSION:
	//    1
	//
	// COMMANDS:
	//    version, v  print the version
	//    help, h     Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    -i int         An integer (default: 22)
	//    --help, -h     show help
	//    --version, -v  print the version
}

func ExampleConfiguration_basic() {
	var (
		cfg configExample
		cmd = &cli.Command{
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("config is", cfg)
				return nil
			},
			Name:    "configbasic",
			Version: "1",
		}
		loaders = []Loader{
			{
				Provider: func(s string) koanf.Provider {
					return file.Provider(s)
				},
				Parser: yaml.Parser(),
				Match: func(s string) bool {
					return strings.HasSuffix(s, ".yml")
				},
			},
		}
	)
	os.Args = []string{"configbasic", "--config", "testdata/test.yml"}
	Run(
		context.Background(),
		cmd,
		Configuration(
			&cfg,
			loaders,
		),
		NoDefaultFlags(),
	)
	// Output:
	// config is {33}
}

func ExampleConfigFlags_basicHelp() {
	// Help for flags bound to fields in a configuration struct. This
	// shows the flag "--i" bound to cfg.I
	//
	var (
		cfg configExample
		cmd = &cli.Command{
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("config is", cfg)
				return nil
			},
			Name:    "basichelp",
			Version: "1",
		}
		loaders = []Loader{
			{
				Provider: func(s string) koanf.Provider {
					return file.Provider(s)
				},
				Parser: yaml.Parser(),
				Match: func(s string) bool {
					return strings.HasSuffix(s, ".yml")
				},
			},
		}
	)
	os.Args = []string{"basichelp", "--help"}
	Run(
		context.Background(),
		cmd,
		Configuration(
			&cfg,
			loaders,
		),
		ConfigFlags(
			[]Configurator{&cfg},
			cmd,
		),
		NoDefaultFlags(),
	)
	// Output:
	// NAME:
	//    basichelp - A new cli application
	//
	// USAGE:
	//    basichelp [global options] [command [command options]]
	//
	// VERSION:
	//    1
	//
	// COMMANDS:
	//    version, v  print the version
	//    help, h     Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    -i value                                                         (default: 0) [$I]
	//    --config string, --cfg string [ --config string, --cfg string ]  comma-separated list of path(s) to configuration file(s)
	//    --help, -h                                                       show help
	//    --version, -v                                                    print the version
}

func ExampleConfigFlags_basicFlagOverride() {
	// Flags bound to fields in a configuration struct. The flag value
	// overrides the field's value obtained from reading the  YAML
	// configuration file
	//
	// test.yml simply contains:
	// i: 33
	//
	var (
		cfg configExample
		cmd = &cli.Command{
			Action: func(ctx context.Context, cmd *cli.Command) error {
				fmt.Println("config is", cfg)
				return nil
			},
			Name:    "basicflagoverride",
			Version: "1",
		}
		loaders = []Loader{
			{
				Provider: func(s string) koanf.Provider {
					return file.Provider(s)
				},
				Parser: yaml.Parser(),
				Match: func(s string) bool {
					return strings.HasSuffix(s, ".yml")
				},
			},
		}
	)
	os.Args = []string{"basicflagoverride", "-i", "77", "--config", "testdata/test.yml"}
	Run(
		context.Background(),
		cmd,
		Configuration(
			&cfg,
			loaders,
		),
		ConfigFlags(
			[]Configurator{&cfg},
			cmd,
		),
		NoDefaultFlags(),
	)
	// Output:
	// config is {77}
}
