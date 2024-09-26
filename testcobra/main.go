package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bruceesmith/echidna"
	"github.com/bruceesmith/echidna/cobra"
	"github.com/bruceesmith/echidna/logger"
	snake "github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type config struct {
	A int
}

var a = config{3}

func execute(cmd *snake.Command, args []string) error {
	fmt.Println("RunE.execute")
	return nil
}

func main() {
	err := cobra.Init(
		"fred",
		"0.1",
		&a,
		"short description",
		"fred [options]",
		func(cf *config) error {
			return nil
		},
		execute,
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		return
	}
	err = cobra.Run()
	if err != nil {
		if errors.Is(err, echidna.ErrConfigOK) {
			fmt.Fprintf(os.Stderr, "configuration is OK\n")
			return
		} else if !errors.Is(err, pflag.ErrHelp) && !errors.Is(err, echidna.ErrVersion) {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
			return
		} else {
			return
		}
	}
	fmt.Println(a)
	logger.TraceID("one", "trace one")
	logger.TraceID("two", "trace two")
}
