package viper

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
func Init[E any](prog, ver string, config *E, validate func(cf *E) error, flagfuncs ...func()) (err error) {
	echidna.Program = prog
	echidna.Version = ver
	fs := pflag.NewFlagSet("vipe", pflag.ContinueOnError)

	// Define standard command line flags
	var cfg = fs.String("cfg", "", "path to configuration file")
	var checkConfig = fs.Bool("checkcfg", false, "check the configuration and then exit")
	var showHelp = fs.Bool("help", false, "print help and then exit")
	var logging logger.LogLevel
	fs.Var(&logging, "log", "logging level (slog values plus LevelTrace)")
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

	// Read the configuration file
	v := viper.New()
	if len(*cfg) == 0 {
		v.SetConfigName(echidna.Program + ".yml")
		v.SetConfigType("yml")
		v.AddConfigPath(".")
	} else {
		v.SetConfigFile(*cfg)
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
