// Copyright © 2024 Bruce Smith <bruceesmith@gmait.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

package vpr

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/bruceesmith/echidna"
	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/echidna/version"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Init initialises the echidna package for the use case where the program
// does not provide and "comamnds" on the command line and it uses the Viper
// package to handle the configuration
func Init[E any](prog, ver string, config *E, envPrefix string, validate func(cf *E) error, flagfuncs ...func()) (err error) {
	echidna.Program = prog
	echidna.Version = ver
	pflag.CommandLine = pflag.NewFlagSet(os.Args[0], pflag.ContinueOnError)

	// Define standard command line flags
	var (
		cfg          string
		checkConfig  bool
		jason        bool
		logging      logger.LogLevel
		showHelp     bool
		traces       logger.Traces = make(logger.Traces, 0)
		printVersion bool
		verbose      bool
	)

	pflag.StringVar(&cfg, "cfg", "", "path to configuration file")
	pflag.BoolVar(&checkConfig, "checkcfg", false, "check the configuration and then exit")
	pflag.BoolVar(&showHelp, "help", false, "print help and then exit")
	pflag.BoolVar(&jason, "json", false, "output should be in JSON format")
	viper.BindPFlag("json", pflag.Lookup("json"))
	pflag.Var(&logging, "log", "logging level (slog values plus LevelTrace)")
	pflag.Var(&traces, "trace", `comma-separated list of trace areas ["all" for every possible area]`)
	pflag.BoolVar(&printVersion, "version", false, "print version details and then exit")
	pflag.BoolVar(&verbose, "verbose", false, "output should be verbose")
	viper.BindPFlag("verbose", pflag.Lookup("verbose"))

	// Define any flags using the pflags package, bind any of them to Viper keys, etc.
	for _, f := range flagfuncs {
		f()
	}

	// Parse all the flags
	err = pflag.CommandLine.Parse(os.Args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
		return
	}

	// Handle "show help"
	if showHelp {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		pflag.PrintDefaults()
		return pflag.ErrHelp
	}

	// Handle "print version"
	if printVersion {
		version.Version()
		return echidna.ErrVersion
	}

	// Prepare to handle environment variables
	if len(envPrefix) > 0 {
		viper.SetEnvPrefix(envPrefix)
	}
	viper.AutomaticEnv()

	// Read the configuration file
	if len(cfg) == 0 {
		viper.SetConfigName(echidna.Program + ".yml")
		viper.SetConfigType("yml")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigFile(cfg)
	}
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	// Set the logging level
	logger.SetLevel(slog.Level(logging))

	// Register areas to be traced, if any
	if len(traces) != 0 {
		logger.SetTraceIds(traces...)
	}

	// Extract the configuration into the provided struct
	err = viper.Unmarshal(config)
	if err != nil {
		return
	}

	// If "checkcfg" is given, validate that the configuration makes sense
	if checkConfig {
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
