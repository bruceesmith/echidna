// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

package custom

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	"github.com/bruceesmith/echidna"
	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/echidna/version"

	"github.com/spf13/pflag"
)

// Init initialises the echidna package for the use case where the program
// does not provide and "comamnds" on the command line and it uses a custom
// package to handle the configuration
func Init[E any](prog, ver string, config *E, unmarshal func([]byte, *E) error, validate func(cf *E) error, flagfuncs ...func()) (err error) {
	echidna.Program = prog
	echidna.Version = ver
	fs := pflag.NewFlagSet("custom", pflag.ContinueOnError)

	// Define standard command ine flags
	var cfg = fs.String("cfg", "", "path to configuration file")
	var checkConfig = fs.Bool("checkcfg", false, "check the configuration and then exit")
	var showHelp = fs.Bool("help", false, "print help and then exit")
	var logging logger.LogLevel
	fs.Var(&logging, "log", "logging level (slog values plus LevelTrace)")
	var traces logger.Traces
	fs.Var(&traces, "trace", `comma-separated list of trace areas ["all" for every possible area]`)
	var printVersion = fs.Bool("version", false, "print version details and then exit")

	// Define any flags using the pflags package, bind any of them to Viper keys, etc.
	for _, f := range flagfuncs {
		f()
	}

	// Parse all the flags
	err = fs.Parse(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()
		return
	}

	// Handle "show help"
	if *showHelp {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		fs.PrintDefaults()
		return pflag.ErrHelp
	}

	// Handle "print version"
	if *printVersion {
		version.Version()
		return echidna.ErrVersion
	}

	// Set the logging level
	logger.SetLevel(slog.Level(logging))

	// Register areas to be traced, if any
	if len(traces) != 0 {
		logger.SetTraceIds(traces...)
	}

	// Read the config file and unmarshal it into the configuration struct
	if len(*cfg) > 0 {
		if unmarshal == nil {
			return errors.New("cannot unmarshal configuration, no unmarshal function provided")
		}
		file, err := os.Open(*cfg)
		if err != nil {
			return err
		}
		defer file.Close()
		bites, err := io.ReadAll(file)
		if err != nil {
			return err
		}
		err = unmarshal(bites, config)
		return err
	}

	// If "checkcfg" is given, validate that the configuration makes sense
	if *checkConfig {
		if validate != nil {
			err = validate(config)
			if err != nil {
				return
			}
			return echidna.ErrConfigOK
		}
		return errors.New("checkcfg requested but no validator function provided")
	}
	return
}
