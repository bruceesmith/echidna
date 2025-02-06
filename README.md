[![Go Reference][goreference_badge]][goreference_link]
[![Go Report Card][goreportcard_badge]][goreportcard_link]
 
# <div align="center">Echidna Library for Go</div>
 
![background image](echidna.png)
 
<!-- Code generated by gomarkdoc. DO NOT EDIT -->

# echidna

```go
import "github.com/bruceesmith/echidna"
```

Package echidna builds upon the Github packages [knadh/koanf](<https://github.com/knadh/koanf>), [urfave/cli/v3](<https://github.com/urfave/cli>), [bruceesmith/sflags](<https://github.com/urfave/sflags>) to make it extremely simple to use the features of these excellent packages in concert.

Every program using echidna will expose a standard set of command\-line flags \(\-\-json, \-\-log, \-\-trace, \-\-verbose\) in addition to the standard flags provided by urfave/cli/v3 \(\-\-help and \-\-version\).

If a configuration struct is provided to [Run](<#Run>) function by [Configuration](<#Configuration>), then a further command\-line flag \(\-\-config\) is added to provide the source\(s\) of values for fields in the struct.

Command\-line flags bound to fields in the configuration are created by providing [ConfigFlags](<#ConfigFlags>) to [Run](<#Run>). These flags can be bound either to the root command or to one or more child commands.

<details><summary>Example (Action)</summary>
<p>



```go
// Include an Action function
(&cli.Command{
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("hello")
		return nil
	},
	Name:    "action",
	Version: "1",
}).Run(context.Background(), []string{"action"})
// Output:
// hello
```

#### Output

```
hello
```

</p>
</details>

<details><summary>Example (Basic)</summary>
<p>



```go
// The most basic example of urfave/cli/v3
(&cli.Command{Name: "basic"}).Run(context.Background(), os.Args)
// Output:
// NAME:
//    basic - A new cli application
//
// USAGE:
//    basic [global options]
//
// GLOBAL OPTIONS:
//    --help, -h  show help
```

#### Output

```
NAME:
   basic - A new cli application

USAGE:
   basic [global options]

GLOBAL OPTIONS:
   --help, -h  show help
```

</p>
</details>

<details><summary>Example (Flag1)</summary>
<p>



```go
// Include a custom flag
(&cli.Command{
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("hello")
		return nil
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "i",
			Usage: "An integer",
		},
	},
	Name:    "action",
	Version: "1",
}).Run(context.Background(), []string{"action"})
// Output:
// hello
```

#### Output

```
hello
```

</p>
</details>

<details><summary>Example (Flag2)</summary>
<p>



```go
// Include a custom flag
(&cli.Command{
	Action: func(ctx context.Context, cmd *cli.Command) error {
		fmt.Println("hello")
		return nil
	},
	Flags: []cli.Flag{
		&cli.IntFlag{
			Name:  "i",
			Usage: "An integer",
			Value: 22,
		},
	},
	Name:    "action",
	Version: "1",
}).Run(context.Background(), []string{"action", "--help"})
// Output:
// NAME:
//    action - A new cli application
//
// USAGE:
//    action [global options]
//
// VERSION:
//    1
//
// GLOBAL OPTIONS:
//    -i value       An integer (default: 22)
//    --help, -h     show help
//    --version, -v  print the version
```

#### Output

```
NAME:
   action - A new cli application

USAGE:
   action [global options]

VERSION:
   1

GLOBAL OPTIONS:
   -i value       An integer (default: 22)
   --help, -h     show help
   --version, -v  print the version
```

</p>
</details>

<details><summary>Example (Version)</summary>
<p>



```go
// Include a Version field in the Command
(&cli.Command{
	Name:    "version",
	Version: "1",
}).Run(context.Background(), os.Args)
// Output:
// NAME:
//    version - A new cli application
//
// USAGE:
//    version [global options]
//
// VERSION:
//    1
//
// GLOBAL OPTIONS:
//    --help, -h     show help
//    --version, -v  print the version
```

#### Output

```
NAME:
   version - A new cli application

USAGE:
   version [global options]

VERSION:
   1

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

</p>
</details>

## Index

- [Variables](<#variables>)
- [func Run\(ctx context.Context, command \*cli.Command, options ...Option\)](<#Run>)
- [type Configurator](<#Configurator>)
- [type FlagOption](<#FlagOption>)
  - [func DescTag\(tag string\) FlagOption](<#DescTag>)
  - [func EnvDivider\(divider string\) FlagOption](<#EnvDivider>)
  - [func EnvPrefix\(prefix string\) FlagOption](<#EnvPrefix>)
  - [func FlagDivider\(divider string\) FlagOption](<#FlagDivider>)
  - [func FlagTag\(tag string\) FlagOption](<#FlagTag>)
  - [func Flatten\(flatten bool\) FlagOption](<#Flatten>)
  - [func Prefix\(prefix string\) FlagOption](<#Prefix>)
  - [func Validator\(val sflags.ValidateFunc\) FlagOption](<#Validator>)
- [type Option](<#Option>)
  - [func ConfigFlags\(configs \[\]Configurator, command \*cli.Command, ops ...FlagOption\) Option](<#ConfigFlags>)
  - [func Configuration\(config Configurator\) Option](<#Configuration>)
  - [func NoDefaultFlags\(\) Option](<#NoDefaultFlags>)
  - [func NoJSON\(\) Option](<#NoJSON>)
  - [func NoLog\(\) Option](<#NoLog>)
  - [func NoTrace\(\) Option](<#NoTrace>)
  - [func NoVerbose\(\) Option](<#NoVerbose>)


## Variables

<a name="BuildDate"></a>

```go
var (
    // BuildDate is the timestamp for when this program was compiled
    BuildDate string = `Filled in during the build`
)
```

<a name="Run"></a>
## func [Run](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L413>)

```go
func Run(ctx context.Context, command *cli.Command, options ...Option)
```

Run is the primary external function of this library. It augments the cli.Command with default command\-line flags, hooks in handling for processing a configuration, runs the appropriate Action, calls the terminator to wait for goroutine cleanup

<a name="Configurator"></a>
## type [Configurator](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L47-L49>)

Configurator is the interface for a configuration struct

```go
type Configurator interface {
    Validate() error
}
```

<a name="FlagOption"></a>
## type [FlagOption](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L43>)

FlagOption is a functional option parameter for the ConfigFlags function

```go
type FlagOption func()
```

<a name="DescTag"></a>
### func [DescTag](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L66>)

```go
func DescTag(tag string) FlagOption
```

DescTag sets the struct tag where usage text is configured

<a name="EnvDivider"></a>
### func [EnvDivider](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L74>)

```go
func EnvDivider(divider string) FlagOption
```

EnvDivider is the character in between parts of an environment variable bound to a command line struct\-bound flag

<a name="EnvPrefix"></a>
### func [EnvPrefix](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L82>)

```go
func EnvPrefix(prefix string) FlagOption
```

EnvPrefix is an optional prefix for an environment variable bound to a command line struct\-bound flag

<a name="FlagDivider"></a>
### func [FlagDivider](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L90>)

```go
func FlagDivider(divider string) FlagOption
```

FlagDivider is the character in between parts of a struct\-bound command line flag's name

<a name="FlagTag"></a>
### func [FlagTag](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L98>)

```go
func FlagTag(tag string) FlagOption
```

FlagTag is the struct tag used to configure struct\-bound command line flags

<a name="Flatten"></a>
### func [Flatten](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L106>)

```go
func Flatten(flatten bool) FlagOption
```

Flatten determines the name of a command line flag bound to a an anonymous struct field

<a name="Prefix"></a>
### func [Prefix](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L114>)

```go
func Prefix(prefix string) FlagOption
```

Prefix is an optional prefix for the names of all struct\-bound command line flags

<a name="Validator"></a>
### func [Validator](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L122>)

```go
func Validator(val sflags.ValidateFunc) FlagOption
```

Validator is an optional function that will be called to to validate each struct\-field bound flag

<a name="Option"></a>
## type [Option](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L59>)

Option is a functional parameter for Run\(\)

```go
type Option func() error
```

<a name="ConfigFlags"></a>
### func [ConfigFlags](<https://github.com/bruceesmith/echidna/blob/main/config_flags.go#L131>)

```go
func ConfigFlags(configs []Configurator, command *cli.Command, ops ...FlagOption) Option
```

ConfigFlags creates and configures [Run](<#Run>) to have command line flags bound to the fields of one or more parts of a configuration struct

<a name="Configuration"></a>
### func [Configuration](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L237>)

```go
func Configuration(config Configurator) Option
```

Configuration is an Option helper to define a configuration structure that will be populated from the sources given on a \-\-config command\-line flag

<a name="NoDefaultFlags"></a>
### func [NoDefaultFlags](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L455>)

```go
func NoDefaultFlags() Option
```

NoDefaultFlags is a convenience function which is equivalent to calling all of NoJSON, NoLog, NoTrace, and NoVerbose

<a name="NoJSON"></a>
### func [NoJSON](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L464>)

```go
func NoJSON() Option
```



<a name="NoLog"></a>
### func [NoLog](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L471>)

```go
func NoLog() Option
```



<a name="NoTrace"></a>
### func [NoTrace](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L478>)

```go
func NoTrace() Option
```



<a name="NoVerbose"></a>
### func [NoVerbose](<https://github.com/bruceesmith/echidna/blob/main/echidna.go#L485>)

```go
func NoVerbose() Option
```



Generated by [gomarkdoc](<https://github.com/princjef/gomarkdoc>)
 
[goreference_badge]: https://pkg.go.dev/badge/github.com/bruceesmith/echidna/v3.svg
[goreference_link]: https://pkg.go.dev/github.com/bruceesmith/echidna
[goreportcard_badge]: https://goreportcard.com/badge/github.com/bruceesmith/echidna
[goreportcard_link]: https://goreportcard.com/report/github.com/bruceesmith/echidna
