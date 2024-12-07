// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

/*
Package version supports either printing or return of information concerning
the compiled CLI or daemon.

[Print] writes a text version of [Info] to [os.Stdout] by default.
- if jason is true, the output is indented JSON instead of text
- if an optional [io.Writer] is supplied the output is written using that [io.Writer].

Alternatively, [Version] returns an [Info] struct

Examples:

	// Print version information to stdout as human readable text
	version.Print(false)

	// Print version information to stdout as indented JSON
	version.Print(true)

	// Print version information into a bytes.Buffer
	b := bytes.NewBufferString("")
	version.Print(false, b)
	fmt.Println(b.String())

	// Return the version information as an Info struct
	info := version.Version()
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", v)
*/
package version

import (
	"encoding/json"
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
	Program    string `json:"program,omitempty"`
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
	if len(vi.Version) == 0 {
		mod := &bi.Main
		if mod.Replace != nil {
			mod = mod.Replace
		}
		vi.Version = mod.Version
	}
	return
}

// Print prints detailed version information
func Print(jason bool, w ...io.Writer) {
	var writer io.Writer = os.Stdout
	if len(w) > 0 {
		writer = w[0]
	}
	vi := makeVersion()
	if jason {
		bites, err := json.MarshalIndent(vi, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, `{"error": "unable to marshal version information: %v"}\n`, err)
		} else {
			fmt.Fprintln(writer, string(bites))
		}
	} else {
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
}

// Version returns detailed version information
func Version() Info {
	return makeVersion()
}
