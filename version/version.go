// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

/*
Package version supports either printing or return of information concerning
the compiled CLI or daemon.

[Print] writes a text version of [Info] to [os.Stderr] by default. If an optional [io.Writer]
is supplied the output is written using that [io.Writer].

Alternatively, [Version] returns an [Info] struct

Examples:

	// Print version information to stderr
	version.Print()

	// Print version information into a bytes.Buffer
	b := bytes.NewBufferString("")
	version.Print(b)
	fmt.Println(b.String())

	// Return the version information as an Info struct
	info := version.Version
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", v)
*/
package version

import (
	"fmt"
	"io"
	"os"
	"runtime/debug"

	"github.com/bruceesmith/echidna"
)

// Info is facts about how and when the program was compiled
// including Git details
type Info struct {
	BuildDate  string `json:"build_date,omitempty"`
	Commit     string `json:"vcs_commit,omitempty"`
	CommitDate string `json:"vcs_date,omitempty"`
	Go         string `json:"go,omitempty"`
	Program    string `json:"progra,omitempty"`
	Version    string `json:"version,omitempty"`

	Error string `json:"error,omitempty"`
}

// makeVersion assembles the version information
func makeVersion() (vi Info) {
	vi.BuildDate = echidna.BuildDate
	bi, ok := debug.ReadBuildInfo()
	if ok {
		vi.Go = bi.GoVersion
		for _, bs := range bi.Settings {
			if bs.Key == "vcs.revision" {
				vi.Commit = bs.Value
			}
			if bs.Key == "vcs.time" {
				vi.CommitDate = bs.Value
			}
		}
	}
	vi.Program = echidna.Program
	vi.Version = echidna.Version
	return
}

// Print prints detailed version information
func Print(w ...io.Writer) {
	var writer io.Writer = os.Stdout
	if len(w) > 0 {
		writer = w[0]
	}
	vi := makeVersion()
	fmt.Fprintf(writer, "%s %s\n", vi.Program, vi.Version)
	fmt.Fprintf(writer, "Compiled with Go %s\n", vi.Go)
	fmt.Fprintf(writer, "Built at %s\n", vi.BuildDate)
	if len(vi.Commit) > 0 {
		fmt.Fprintf(writer, "VCS commit %s\n", vi.Commit)
	}
	if len(vi.CommitDate) > 0 {
		fmt.Fprintf(writer, "VCS commit time %s\n", vi.CommitDate)
	}
}

// Version returns detailed version information
func Version() Info {
	return makeVersion()
}
