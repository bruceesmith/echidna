module github.com/bruceesmith/echidna

go 1.24.0

require (
	github.com/bruceesmith/logger v1.2.0
	github.com/bruceesmith/sflags v0.3.2
	github.com/bruceesmith/terminator v1.0.0
	github.com/deckarep/golang-set/v2 v2.8.0
	github.com/jinzhu/copier v0.4.0
	github.com/knadh/koanf v1.5.0
	github.com/stretchr/testify v1.10.0
	github.com/urfave/cli/v3 v3.0.0-beta1

)

require (
	github.com/BurntSushi/toml v1.5.0 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.6 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fatih/structs v1.1.0 // indirect
	github.com/fsnotify/fsnotify v1.8.0 // indirect
	github.com/gojp/goreportcard v0.0.0-20250313055930-6c0fa925565c // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.27.6 // indirect
	github.com/xrash/smetrics v0.0.0-20240521201337-686a1a2994c1 // indirect
	golang.org/x/exp/typeparams v0.0.0-20250305212735-054e65f0b394 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.12.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/telemetry v0.0.0-20250310203348-fdfaad844314 // indirect
	golang.org/x/tools v0.31.0 // indirect
	golang.org/x/vuln v1.1.4 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	honnef.co/go/tools v0.6.1 // indirect
)

tool (
	github.com/gojp/goreportcard/cmd/goreportcard-cli
	golang.org/x/vuln/cmd/govulncheck
	honnef.co/go/tools/cmd/staticcheck
)
