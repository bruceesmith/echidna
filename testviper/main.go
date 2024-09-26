package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bruceesmith/echidna"
	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/echidna/vpr"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type config struct {
	A int
}

var a = config{2}

var toggle string

func main() {
	flags := func() {
		viper.BindPFlag("log", pflag.Lookup("log"))
		pflag.StringVar(&toggle, "toggle", "", "test flag for toggle")
	}
	err := vpr.Init(
		"fred",
		"0.1",
		&a,
		"",
		func(cf *config) error {
			return nil
		},
		flags,
	)
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
	fmt.Println(viper.GetString("log"))
	fmt.Println(viper.GetBool("json"))
	fmt.Println(viper.GetBool("verbose"))
}
