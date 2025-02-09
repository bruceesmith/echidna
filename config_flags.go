// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package echidna

import (
	"fmt"
	"reflect"
	"slices"

	"github.com/bruceesmith/sflags"
	"github.com/bruceesmith/sflags/gen/gcli"
	"github.com/jinzhu/copier"
	"github.com/urfave/cli/v3"
)

// binder is the heart of the flag-binding mechanism:
//   - a command-line flag can be bound to a struct field by use of tags
//     on the field, and by providing all or part of a struct as an
//     argument on a [ConfigFlags] call
//   - when a [urfave/cli/v3/Run] call is executed, the command-line
//     arguments are parsed into the associated struct fields
//   - a copy of the updated struct is made, thereby preserving the
//     values of any struct-bound flag that was provided on the command line
//   - the command's Before function is executed, and it loads the
//     configuration from sources specified on the command-line, over-writing
//     the struct-bound values in the configuration
//   - the values of struct-bound flags from the copy of the configuration
//     are pushed back into the configuration
//   - execution of the command's Action proceeds
//   - as a result, values for struct-bound fields are set in this order
//     1. default value in the configuration struct
//     2. value set in one of the configuration sources loaded by [knadh/koanf]
//     3. environment variable configured by [bruceesmith/sflags]
//     4. flag value from the command line
type binder struct {
	clone        Configurator
	configFields map[string]reflect.Value
	cloneFields  map[string]reflect.Value
}

// FlagOption is a functional option parameter for the ConfigFlags function
type FlagOption func()

// options are settings used to configure and then handle command line
// flags bound to configuration struct fields
type options struct {
	descTag, divider, envDivider, envPrefix, tag, prefix string
	flatten                                              bool
	validator                                            sflags.ValidateFunc
}

var (
	opts = options{
		descTag:    "desc",
		divider:    "-",
		envDivider: "_",
		envPrefix:  "",
		tag:        "flag",
		flatten:    false,
		prefix:     "",
	} // Default option settings
)

// DescTag sets the struct tag where usage text is configured
func DescTag(tag string) FlagOption {
	return func() {
		opts.descTag = tag
	}
}

// EnvDivider is the character in between parts of an environment
// variable bound to a command line struct-bound flag
func EnvDivider(divider string) FlagOption {
	return func() {
		opts.envDivider = divider
	}
}

// EnvPrefix is an optional prefix for an environment
// variable bound to a command line struct-bound flag
func EnvPrefix(prefix string) FlagOption {
	return func() {
		opts.envPrefix = prefix
	}
}

// FlagDivider is the character in between parts of a
// struct-bound command line flag's name
func FlagDivider(divider string) FlagOption {
	return func() {
		opts.divider = divider
	}
}

// FlagTag is the struct tag used to configure struct-bound
// command line flags
func FlagTag(tag string) FlagOption {
	return func() {
		opts.tag = tag
	}
}

// Flatten determines the name of a command line flag
// bound to a an anonymous struct field
func Flatten(flatten bool) FlagOption {
	return func() {
		opts.flatten = flatten
	}
}

// Prefix is an optional prefix for the names of all
// struct-bound command line flags
func Prefix(prefix string) FlagOption {
	return func() {
		opts.prefix = prefix
	}
}

// Validator is an optional function that will be
// called to to validate each struct-field bound flag
func Validator(val sflags.ValidateFunc) FlagOption {
	return func() {
		opts.validator = val
	}
}

// ConfigFlags creates and configures [Run] to have
// command line flags bound to the fields of one or more
// parts of a configuration struct
func ConfigFlags(configs []Configurator, command *cli.Command, ops ...FlagOption) Option {
	return func() error {
		if len(configs) == 0 {
			return fmt.Errorf("ConfigFlags called with zero configuration structs")
		}
		for _, opt := range ops {
			opt()
		}
		flags := make([]cli.Flag, 0)
		for i, cfg := range configs {
			flgs := make([]cli.Flag, 0)
			err := gcli.ParseToV3(
				cfg,
				&flgs,
				sflags.DescTag(opts.descTag),
				sflags.EnvDivider(opts.envDivider),
				sflags.EnvPrefix(opts.envPrefix),
				sflags.FlagDivider(opts.divider),
				sflags.FlagTag(opts.tag),
				sflags.Flatten(opts.flatten),
				sflags.Prefix(opts.prefix),
			)
			if err != nil {
				return fmt.Errorf("ConfigFlags() failed to create []cli.Flags for configuration %d: [%w]", i, err)
			}
			flags = slices.Grow(flags, len(flags)+len(flgs))
			flags = append(flags, flgs...)
		}
		if len(flags) > 0 {
			addFlags(command, flags)
		}
		return nil
	}
}

// applyFlagOverrides applies the values of struct-bound command
// line flags to the associated struct fields. It only does so
// for such flags that are present on the command line
func applyFlagOverrides(names []string, b binder) {
	for _, f := range names {
		var ok bool
		var in, out reflect.Value
		in, ok = b.cloneFields[f]
		if !ok {
			continue
		}
		out, ok = b.configFields[f]
		if !ok {
			continue
		}
		out.Set(in)
	}

}

// clone safely creates an exact copy of a configuration struct
func clone(in Configurator) Configurator {
	out := reflect.New(reflect.ValueOf(in).Elem().Type()).Interface().(Configurator)
	copier.Copy(out, in)
	return out
}

// fieldMap returns a structure that associates flattened field
// names with struct fields
func fieldMap(v Configurator) (map[string]reflect.Value, error) {
	if v == nil {
		return nil, fmt.Errorf("fieldMap requires a non-nil pointer")
	}
	if reflect.ValueOf(v).Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("fieldMap requires a pointer to a struct")
	}
	result := make(map[string]reflect.Value)
	mapFields(reflect.ValueOf(v).Elem(), "", &result)
	return result, nil
}

// mapFields is a recursive function which traverses a struct
// to build a field map
func mapFields(v reflect.Value, prefix string, result *map[string]reflect.Value) {
	if len(prefix) != 0 {
		prefix += opts.divider
	}
	value := v
	tipe := value.Type()
	if tipe.Kind() == reflect.Struct {
		for i := 0; i < tipe.NumField(); i++ {
			field := tipe.Field(i)
			if field.Type.Kind() != reflect.Struct {
				(*result)[prefix+parseField(field, opts.prefix)] = value.Field(i)
			} else {
				mapFields(value.Field(i), prefix+camelToFlag(field.Name, opts.divider), result)
			}
		}
	}
}

// newFlagBinder creates a binder, the basis for associating and setting
// configuration struct fields from command line flags
func newFlagBinder(cfg Configurator) (b binder, err error) {
	b.clone = clone(cfg)
	b.configFields, err = fieldMap(cfg)
	if err != nil {
		return b, fmt.Errorf("cannot build a field map for the configuration: [%w]", err)
	}
	b.cloneFields, err = fieldMap(b.clone)
	if err != nil {
		return b, fmt.Errorf("cannot build a field map for the configuration clone: [%w]", err)
	}
	return b, nil
}
