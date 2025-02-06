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
	// The most basic example of urfave/cli/v3
	(&cli.Command{Name: "basic"}).Run(context.Background(), os.Args)
	// Output:
	// NAME:
	//    basic - A new cli application
	//
	// USAGE:
	//    basic [global options]
	//
	// GLOBAL OPTIONS:
	//    --help, -h  show help
}

func Example_version() {
	// Include a Version field in the Command
	(&cli.Command{
		Name:    "version",
		Version: "1",
	}).Run(context.Background(), os.Args)
	// Output:
	// NAME:
	//    version - A new cli application
	//
	// USAGE:
	//    version [global options]
	//
	// VERSION:
	//    1
	//
	// GLOBAL OPTIONS:
	//    --help, -h     show help
	//    --version, -v  print the version
}

func Example_action() {
	// Include an Action function
	(&cli.Command{
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("hello")
			return nil
		},
		Name:    "action",
		Version: "1",
	}).Run(context.Background(), []string{"action"})
	// Output:
	// hello
}

func Example_flag1() {
	// Include a custom flag
	(&cli.Command{
		Action: func(ctx context.Context, cmd *cli.Command) error {
			fmt.Println("hello")
			return nil
		},
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  "i",
				Usage: "An integer",
			},
		},
		Name:    "action",
		Version: "1",
	}).Run(context.Background(), []string{"action"})
	// Output:
	// hello
}

func Example_flag2() {
	// Include a custom flag
	(&cli.Command{
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
		Name:    "action",
		Version: "1",
	}).Run(context.Background(), []string{"action", "--help"})
	// Output:
	// NAME:
	//    action - A new cli application
	//
	// USAGE:
	//    action [global options]
	//
	// VERSION:
	//    1
	//
	// GLOBAL OPTIONS:
	//    -i value       An integer (default: 22)
	//    --help, -h     show help
	//    --version, -v  print the version
}
