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
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"reflect"
	"runtime/debug"
	"strings"

	"github.com/bruceesmith/echidna/logger"
	"github.com/bruceesmith/echidna/terminator"
	"github.com/knadh/koanf"
	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/urfave/cli/v3"
)

// Configuration is the interface for a configuration struct
type Configuration interface {
	Validate() error
}

// configLoader is a parameter to koanf.Load()
type configLoader struct {
	Provider koanf.Provider
	Parser   koanf.Parser
	Options  []koanf.Option
}

// Option is a functional parameter for Run()
type Option func(params ...any) error

var (
	// Flag to specifiy the configuration source(s)
	configFlag = &cli.StringSliceFlag{
		Name:    "config",
		Aliases: []string{"cfg"},
		Usage:   "comma-separated list of path(s) to configuration file(s)",
	}
	// Pointer to a configuration struct
	configuration Configuration
	// Default command-line flags beyond --help and --version that
	// are provided natively by Koanf
	standardFlags = []cli.Flag{
		&cli.BoolFlag{
			Name:    "json",
			Aliases: []string{"J"},
			Usage:   "output should be JSON format",
		},
		&logger.LogLevelFlag{
			Name:  "log",
			Usage: "logging level (slog values plus LevelTrace)",
		},
		&cli.StringSliceFlag{
			Name:  "trace",
			Usage: `comma-separated list of trace areas ["all" for every possible area]`,
		},
		&cli.BoolFlag{
			Name:    "verbose",
			Aliases: []string{"V"},
			Usage:   "verbose output",
		},
	}
	// Command to print version information
	version = &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the version",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			printVersion(cmd)
			return nil
		},
	}
)

func init() {
	// Set a custom version printer
	cli.VersionPrinter = printVersion
}

// addCommand adds a subcommand to an existing cli.Command
func addCommand(cmd *cli.Command, command *cli.Command) {
	cmd.Commands = append(cmd.Commands, command)
}

// addFlags adds one or more cli.Flag definitions to a command
func addFlags(cmd *cli.Command, flags []cli.Flag) {
	cmd.Flags = append(cmd.Flags, flags...)
}

// before is executed by cmd.Run() after the command line has been processed
// but prior to executing the Action
func before(ctx context.Context, cmd *cli.Command) (cctx context.Context, err error) {
	// Set up logging
	if err = logging(cmd); err != nil {
		return ctx, fmt.Errorf("command initialisation failed: [%w]", err)
	}
	// Read, parse, validate and store the configuration
	if configuration != nil {
		configs := cmd.StringSlice("config")
		if len(configs) == 0 {
			return ctx, fmt.Errorf("configuration not set, no --config flag provided")
		} else {
			// Build a list of configuration source providers
			var configloaders []configLoader
			configloaders, err = loaders(configs)
			// Read, parse, store and validate the configuration
			err = configure(configuration, configloaders)
		}
	}
	return ctx, err
}

// logging establishes logging according to any relevant command-line flags
func logging(command *cli.Command) error {
	// Change the logging format if JSON was requested
	if command.Bool("json") {
		logger.SetFormat(logger.JSON)
	}
	// Set the logging level
	str, set := flag(command, "log")
	if set {
		var ll logger.LogLevel
		if err := ll.Set(str); err != nil {
			return fmt.Errorf("failed to set logging level: [%w]", err)
		}
		logger.SetLevel(slog.Level(ll))
	}
	// Register areas to be traced, if any
	traces := command.StringSlice("trace")
	if len(traces) != 0 {
		logger.SetTraceIds(traces...)
	}
	return nil
}

// configure reads the configuration from the nominated sources, unmarshals it into
// the provided struct, and finally invokes the configuration validator
func configure(config Configuration, configLoaders []configLoader) (err error) {
	konfigurator := koanf.New(".")
	err = readConfig(konfigurator, configLoaders...)
	if err != nil {
		return fmt.Errorf("failed to load configuration: [%w]", err)
	}

	err = konfigurator.Unmarshal("", config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal configuration: [%w]", err)
	}

	err = config.Validate()
	if err != nil {
		return fmt.Errorf("configuration validation failed: [%w]", err)
	}

	return
}

// flag returns the string value of a command-line flag
// found is true and value is non-empty if the flag has been set
// If a flag is one of the standard cli-defined types (e.g. BoolFlag, StringFlag, and so on)
// then using native methods (cmd.Bool(), cmd.String(), and so on) is preferred to use of flag.
// flag is useful for custom FlagSets
func flag(cmd *cli.Command, name string) (value string, found bool) {
	has := func(names []string, name string) bool {
		for _, n := range names {
			if n == name {
				return true
			}
		}
		return false
	}
	for _, flag := range cmd.Flags {
		if has(flag.Names(), name) {
			if flag.IsSet() {
				return flag.String(), true
			}
		}
	}
	return
}

// readConfig reads the configuration from the nominated sources
func readConfig(k *koanf.Koanf, sources ...configLoader) error {
	var err, result error
	for _, source := range sources {
		err = k.Load(source.Provider, source.Parser, source.Options...)
		if err != nil {
			if result != nil {
				result = fmt.Errorf("%s; %s", result.Error(), err.Error())
			} else {
				result = err
			}
		}
	}
	return result
}

// loaders constructs a configuration loader for each nominated source
func loaders(paths []string) ([]configLoader, error) {
	loaders := make([]configLoader, len(paths), len(paths))
	for i, path := range paths {
		loader := configLoader{
			Provider: file.Provider(path),
		}
		switch {
		case strings.HasSuffix(path, "json"):
			loader.Parser = kjson.Parser()
		case strings.HasSuffix(path, "yml"), strings.HasSuffix(path, "yaml"):
			loader.Parser = yaml.Parser()
		default:
			err := fmt.Errorf("no configuraiton parser defined for %s", path)
			return nil, err
		}
		loaders[i] = loader
	}
	return loaders, nil
}

// printVersion is a custom function to print version information
func printVersion(cmd *cli.Command) {
	type ver struct {
		Name      string `json:"name"`
		Version   string `json:"version"`
		GoVersion string `json:"go_version,omitempty"`
		Commit    string `json:"commit,omitempty"`
	}
	info := ver{
		Name:    cmd.Name,
		Version: cmd.Version,
	}
	if cmd.Bool("verbose") {
		bi, ok := debug.ReadBuildInfo()
		if ok {
			info.GoVersion = bi.GoVersion
			for _, v := range bi.Settings {
				if v.Key == "vcs.revision" {
					info.Commit = v.Value
					break
				}
			}
		}
	}
	if cmd.Bool("json") {
		bites, err := json.Marshal(info)
		if err != nil {
			fmt.Println(`{"error":"` + err.Error() + `"}`)
		} else {
			fmt.Println(string(bites))
		}
	} else {
		fmt.Println(info.Name, info.Version)
		if cmd.Bool("verbose") {
			fmt.Println("Compiled with Go version", info.GoVersion)
			if info.Commit != "" {
				fmt.Println("Git commit", info.Commit)
			}
		}
	}

}

// Run is the primary external function of this library. It augments the
// cli.Command with default command-line flags, hooks in handling for
// processing a configuration, runs the appropriate Action, calls the
// terminator to wait for goroutine cleanup
func Run(ctx context.Context, command *cli.Command, options ...Option) {
	var err error
	addFlags(command, standardFlags)
	addCommand(command, version)
	for _, opt := range options {
		err := opt()
		if err != nil {
			logger.Error("Error executing Run() options", "error", err.Error())
			os.Exit(1)
		}
	}
	if configuration != nil {
		addFlags(command, []cli.Flag{configFlag})
	}
	command.Before = before

	err = command.Run(ctx, os.Args)
	terminator.Wait()
	if err != nil && !strings.Contains(err.Error(), "flag provided but not defined") {
		logger.Error("Error performing command", "error", err.Error(), "command", command.FullName())
		os.Exit(1)
	}
}

// WithConfiguration is an Option helper to define a configuration structure
// that will be populated from the sources given on a --config command-line flag
func WithConfiguration(config Configuration) Option {
	return func(_ ...any) error {
		if reflect.TypeOf(config).Kind() != reflect.Pointer {
			return fmt.Errorf("argument to program.WithConfiguration must be a pointer")
		}
		if reflect.TypeOf(config).Elem().Kind() != reflect.Struct {
			return fmt.Errorf("argument to program.WithConfiguration must be a pointer to a struct")
		}
		configuration = config
		return nil
	}
}
