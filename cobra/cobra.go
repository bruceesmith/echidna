// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package cobra

import (
	"errors"
	"log/slog"

	"github.com/bruceesmith/echidna"
	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/echidna/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd *cobra.Command
	initErr error
)

// Run adds all child commands to the root command and sets flags appropriately.
// This is called by the main program. It only needs to happen once to the rootCmd.
func Run() (err error) {
	err = rootCmd.Execute()
	if initErr != nil {
		return initErr
	}
	if err != nil {
		return
	}
	return
}

// Init initialises the echidna package for the use case where the program
// defines "commands" using the [github.com/spf13/cobra] package.
func Init[E any](prog, ver string, config *E, short, use string, validate func(cf *E) error, execute func(*cobra.Command, []string) error) (err error) {
	var (
		cfgPath      string
		checkConfig  bool
		logging      logger.LogLevel = logger.LogLevel(slog.LevelInfo)
		printVersion bool
		traces       logger.Traces
	)
	initialise := func() {
		// Handle "print version"
		if printVersion {
			version.Version()
			initErr = echidna.ErrVersion
			return
		}

		// Set the logging level
		logger.SetLevel(slog.Level(logging))

		// Register areas to be traced, if any
		if len(traces) != 0 {
			logger.SetTraceIds(traces...)
		}

		// Read the configuration file
		v := viper.New()
		if len(cfgPath) == 0 {
			v.SetConfigName(echidna.Program + ".yml")
			v.SetConfigType("yml")
			v.AddConfigPath(".")
		} else {
			v.SetConfigFile(cfgPath)
		}
		err = v.ReadInConfig()
		if err != nil {
			return
		}

		// Extract the configuration into the provided struct
		err = v.Unmarshal(config)
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
				initErr = echidna.ErrConfigOK
			} else {
				initErr = errors.New("checkcfg requested but no validator function provided")
			}
		}
	}
	cobra.OnInitialize(initialise)

	echidna.Program = prog
	echidna.Version = ver
	rootCmd = &cobra.Command{
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// initErr may be set by the OnInitialize function
			if initErr != nil {
				return
			}
			return execute(cmd, args)
		},
		Use: use,
	}

	rootCmd.PersistentFlags().StringVar(&cfgPath, "cfg", "", "path to configuration file")
	rootCmd.Flags().BoolVar(&checkConfig, "checkcfg", false, "check the configuration and then exit")
	rootCmd.PersistentFlags().Var(&logging, "log", "logging level (slog values plus LevelTrace)")
	rootCmd.PersistentFlags().Var(&traces, "trace", `comma-separated list of trace areas ["all" for every possible area]`)
	rootCmd.PersistentFlags().BoolVar(&printVersion, "version", false, "print version details and then exit")
	return
}
