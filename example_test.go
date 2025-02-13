// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package echidna

import (
	"context"
	"fmt"
	"os"

	"github.com/urfave/cli/v3"
)

func Example_basic() {
	// The most basic example
	//
	var cmd = &cli.Command{
		Name: "basic",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
	}
	os.Args = []string{"basic", "h"}
	Run(context.Background(), cmd, NoDefaultFlags())
	// Output:
	// NAME:
	//    basic - A new cli application
	//
	// USAGE:
	//    basic [global options] [command [command options]]
	//
	// COMMANDS:
	//    version, v  print the version
	//    help, h     Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    --help, -h  show help
}

func Example_version() {
	// Include a Version field in the Command
	//
	var cmd = &cli.Command{
		Name: "version",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return nil
		},
		Version: "1",
	}
	os.Args = []string{"basic", "h"}
	Run(context.Background(), cmd, NoDefaultFlags())
	// Output:
	// NAME:
	//    version - A new cli application
	//
	// USAGE:
	//    version [global options] [command [command options]]
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

func Example_action() {
	// Include an Action function
	//
	var cmd = &cli.Command{
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("hello")
			return nil
		},
		Name:    "action",
		Version: "1",
	}
	os.Args = []string{"action"}
	Run(context.Background(), cmd)
	// Output:
	// hello
}

func Example_customflag() {
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
		Name:    "customflag",
		Version: "1",
	}
	os.Args = []string{"action", "-i", "22"}
	Run(context.Background(), cmd, NoDefaultFlags())
	// Output:
	// hello 22
}

func Example_flagwithdefault() {
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
		Name:    "flagwithdefault",
		Version: "1",
	}
	os.Args = []string{"flagwithdefault", "--help"}
	Run(context.Background(), cmd, NoDefaultFlags())
	// Output:
	// NAME:
	//    flagwithdefault - A new cli application
	//
	// USAGE:
	//    flagwithdefault [global options] [command [command options]]
	//
	// VERSION:
	//    1
	//
	// COMMANDS:
	//    version, v  print the version
	//    help, h     Shows a list of commands or help for one command
	//
	// GLOBAL OPTIONS:
	//    -i value       An integer (default: 22)
	//    --help, -h     show help
	//    --version, -v  print the version
}
