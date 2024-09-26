package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bruceesmith/echidna"
	"github.com/bruceesmith/echidna/custom"
	"github.com/bruceesmith/echidna/logger"
	"gopkg.in/yaml.v3"

	"github.com/spf13/pflag"
)

type config struct {
	A int
}

var a = config{2}

func unmarshal(bites []byte, cfg *config) (err error) {
	err = yaml.Unmarshal(bites, cfg)
	return
}

func main() {
	err := custom.Init(
		"fred",
		"0.1",
		&a,
		unmarshal,
		func(cf *config) error {
			return nil
		},
	)
	// prog, ver string, config *E, unmarshal func([]byte, *E) error, validate func(cf *E) error, flagfuncs ...func()unmarshal func([]byte, *E) error, validate func(cf *E) error, flagfuncs ...func()
	if err != nil {
		if errors.Is(err, echidna.ErrConfigOK) {
			fmt.Fprintf(os.Stderr, "configuration is OK")
		} else if !errors.Is(err, pflag.ErrHelp) && !errors.Is(err, echidna.ErrVersion) {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		} else {
			return
		}
	}
	fmt.Println(a)
	logger.TraceID("one", "trace one")
	logger.TraceID("two", "trace two")
}
