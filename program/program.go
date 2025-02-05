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

	"github.com/bruceesmith/echidna/terminator"
	"github.com/bruceesmith/logger"
	"github.com/bruceesmith/set"
	"github.com/knadh/koanf"
	kjson "github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/urfave/cli/v3"
)

// Configurator is the interface for a configuration struct
type Configurator interface {
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

// flagset is used to manage the default flags provided by Run()
type flagset struct {
	all   map[string]cli.Flag
	inuse *set.Set[string]
}

// Delete removes one of the default flags
func (fs flagset) Delete(name string) {
	fs.inuse.Delete(name)
}

// InUse returns a slice of the flags that remain in the standard
// flag set after all Options that remove a flag have been executed
func (fs flagset) InUse() []cli.Flag {
	val := make([]cli.Flag, fs.inuse.Size())
	for k, v := range fs.inuse.Members() {
		val[k] = fs.all[v]
	}
	return val
}

// Len is the number of default flags that will be added to the command line
func (fs flagset) Len() int {
	return fs.inuse.Size()
}

var (
	// Standard flags provided if no Option "No***" functions were called on Run() and
	// if Configuration() was called on Run()
	flags = flagset{
		all: map[string]cli.Flag{
			"config": &cli.StringSliceFlag{
				Name:    "config",
				Aliases: []string{"cfg"},
				Usage:   "comma-separated list of path(s) to configuration file(s)",
			},
			"json": &cli.BoolFlag{
				Name:    "json",
				Aliases: []string{"J"},
				Usage:   "output should be JSON format",
			},
			"log": &logger.LogLevelFlag{
				Name:  "log",
				Usage: "logging level (slog values plus LevelTrace)",
				Value: logger.LogLevel(logger.LevelTrace),
			},
			"trace": &cli.StringSliceFlag{
				Name:  "trace",
				Usage: `comma-separated list of trace areas ["all" for every possible area]`,
			},
			"verbose": &cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"},
				Usage:   "verbose output",
			},
		},
		inuse: set.New(
			"config",
			"json",
			"log",
			"trace",
			"verbose",
		),
	}

	// Pointer to a configuration struct
	configuration Configurator

	// noOsExit is used during testing to avoid calling os.Exit() in Run()
	noOsExit bool

	// Command to print version information
	version = &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "print the version",
		Action: func(ctx context.Context, cmd *cli.Command) error {
			printVersion(cmd.Root())
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
	cmd.Flags = expand(cmd.Flags, len(flags))
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
			return ctx, nil
		} else {
			// The command line has been parsed and values set for any provided flags. If
			// any of the flags were generated from the configuration struct by the [urfave/sflags] package,
			// and any of these mapped flags were provided on the command line, then the associated
			// fields in the configuration struct have been updated from the relevant command line flag(s).
			//
			// The configuration is about to be updated by reading from any configuration sources provided
			// by the flag --config. This would override the flag values that have just been saved. So the
			// configuration is copied at this point. Later, these flag values will be applied again
			// (because flags override values loaded from the configuration sources). Whew ....
			binds, err := newFlagBinder(configuration)
			if err != nil {
				return ctx, fmt.Errorf("configuration handling failed: [%w]", err)
			}

			// Build a list of configuration source providers
			var configloaders []configLoader
			configloaders, err = loaders(configs)
			if err != nil {
				return ctx, fmt.Errorf("config load error: [%w]", err)
			}

			// Read, parse, store the configuration
			err = configure(configuration, configloaders)
			if err != nil {
				return ctx, fmt.Errorf("configuration loading failed: [%w]", err)
			}

			// Update the configuration that has just been loaded with any values that were provided
			// on the command line
			applyFlagOverrides(cmd.FlagNames(), binds)

			// Finally, validate the resulting configuration
			err = configuration.Validate()
			if err != nil {
				return ctx, fmt.Errorf("configuration validation failed: [%w]", err)
			}

		}
	}
	return ctx, err
}

// configure reads the configuration from the nominated sources, unmarshals it into
// the provided struct
func configure(config Configurator, configLoaders []configLoader) (err error) {
	konfigurator := koanf.New(".")
	err = readConfig(konfigurator, configLoaders...)
	if err != nil {
		return fmt.Errorf("failed to load configuration: [%w]", err)
	}

	err = konfigurator.Unmarshal("", config)
	if err != nil {
		return fmt.Errorf("failed to unmarshal configuration: [%w]", err)
	}

	return
}

// Configuration is an Option helper to define a configuration structure
// that will be populated from the sources given on a --config command-line flag
func Configuration(config Configurator) Option {
	return func(_ ...any) error {
		if reflect.TypeOf(config).Kind() != reflect.Pointer {
			return fmt.Errorf("argument to program.Configuration must be a pointer")
		}
		if reflect.TypeOf(config).Elem().Kind() != reflect.Struct {
			return fmt.Errorf("argument to program.Configuration must be a pointer to a struct")
		}
		configuration = config
		return nil
	}
}

// expand grows a slice with either zero or one allocation
func expand[T any](slice []T, size int) (v []T) {
	if cap(slice) < len(slice)+size {
		v = make([]T, len(slice), len(slice)+size)
		copy(v, slice)
		return v
	}
	return slice
}

// flag returns the string value of a command-line flag.
//
// found is true and value is non-empty if the flag has been set.
//
// If a flag is one of the standard cli-defined FlagBase types (e.g. BoolFlag, StringFlag, and so on)
// then using native methods (cmd.Bool(), cmd.String(), and so on) is preferred to use of flag.
// Thus flag is intended to support custom FlagBase types
func flag(cmd *cli.Command, name string) (value any, found bool) {
	has := func(names []string, name string) bool {
		for _, n := range names {
			if n == name {
				return true
			}
		}
		return false
	}

	if has(cmd.FlagNames(), name) {
		// Flag is both defined and provided on the command line
		return cmd.Value(name), true

	}
	for _, flag := range cmd.Flags {
		if has(flag.Names(), name) {
			// Flag is defined but not provided on the command line, so return its default
			type getter interface {
				cli.Flag
				GetValue() string
			}
			fb, ok := flag.(getter)
			if ok {
				return fb.GetValue(), true
			}
			// Flag is defined but is not a FlagBase - it can't be handled in this function
			break
		}
	}
	return value, false
}

// loaders constructs a configuration loader for each nominated source
func loaders(paths []string) ([]configLoader, error) {
	loaders := make([]configLoader, len(paths))
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
			err := fmt.Errorf("no configuration parser defined for %s", path)
			return nil, err
		}
		loaders[i] = loader
	}
	return loaders, nil
}

// logging establishes logging according to any relevant command-line flags
func logging(command *cli.Command) error {
	// Change the logging format if JSON was requested
	if command.Bool("json") {
		logger.SetFormat(logger.JSON)
	}
	// Set the logging level
	value, found := flag(command, "log")
	if found {
		var level logger.LogLevel
		switch lev := value.(type) {
		case string:
			level.Set(lev)
		case logger.LogLevel:
			level = lev
		default:
			return fmt.Errorf("cannot extract a LogLevel from %v type %v", value, reflect.TypeOf(value))
		}
		logger.SetLevel(slog.Level(level))

	}
	// Register areas to be traced, if any
	traces := command.StringSlice("trace")
	if len(traces) != 0 {
		logger.SetTraceIds(traces...)
	}
	return nil
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
			fmt.Fprintln(cmd.Writer, `{"error":"`+err.Error()+`"}`)
		} else {
			fmt.Fprintln(cmd.Writer, string(bites))
		}
	} else {
		fmt.Fprintln(cmd.Writer, info.Name, info.Version)
		if cmd.Bool("verbose") {
			fmt.Fprintln(cmd.Writer, "Compiled with Go version", info.GoVersion)
			if info.Commit != "" {
				fmt.Fprintln(cmd.Writer, "Git commit", info.Commit)
			}
		}
	}

}

// readConfig reads the configuration from the nominated sources
func readConfig(k *koanf.Koanf, sources ...configLoader) error {
	var err, result error
	for _, source := range sources {
		err = k.Load(source.Provider, source.Parser, source.Options...)
		if err != nil {
			if result != nil {
				result = fmt.Errorf("%s: %s", result.Error(), err.Error())
			} else {
				result = err
			}
		}
	}
	return result
}

// Run is the primary external function of this library. It augments the
// cli.Command with default command-line flags, hooks in handling for
// processing a configuration, runs the appropriate Action, calls the
// terminator to wait for goroutine cleanup
func Run(ctx context.Context, command *cli.Command, options ...Option) {
	var err error
	// Apply all the Options
	for _, opt := range options {
		err := opt(command)
		if err != nil {
			logger.Error("Error executing Run() options", "error", err.Error())
			os.Exit(1)
		}
	}
	// No use for a --config flag if Configuration() wasn't used
	if configuration == nil {
		flags.Delete("config")
	}
	// Add on default flags that have not been scrapped
	addFlags(command, flags.InUse())
	// Add a "version" command. Thus seems to be required since we supply
	// our own printVersion function
	addCommand(command, version)
	// Hook in the actions that need to happen after the command line is
	// processed but before the Action code is executed
	command.Before = before
	// Direct logging to the same io.Writers as the command
	if command.Root().Writer != nil {
		logger.RedirectStandard(command.Root().Writer)
	}
	if command.Root().ErrWriter != nil {
		logger.RedirectTrace(command.Root().ErrWriter)
	}

	err = command.Run(ctx, os.Args)
	terminator.Wait()
	if err != nil && !strings.Contains(err.Error(), "flag provided but not defined") {
		logger.Error("Error performing command", "error", err.Error(), "command", command.FullName())
		if !noOsExit {
			os.Exit(1)
		}
	}
}

// NoDefaultFlags is a convenience function which is equivalent to
// calling all of NoJSON, NoLog, NoTrace, and NoVerbose
func NoDefaultFlags() Option {
	return func(_ ...any) error {
		flags.Delete("json")
		flags.Delete("log")
		flags.Delete("trace")
		flags.Delete("verbose")
		return nil
	}
}
func NoJSON() Option {
	return func(_ ...any) error {
		flags.Delete("json")
		return nil
	}
}

func NoLog() Option {
	return func(_ ...any) error {
		flags.Delete("log")
		return nil
	}
}

func NoTrace() Option {
	return func(_ ...any) error {
		flags.Delete("trace")
		return nil
	}
}

func NoVerbose() Option {
	return func(_ ...any) error {
		flags.Delete("verbose")
		return nil
	}
}
