// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

/*
Package echidna provides helpers for building robust Go daemons and CLIs

Type Terminator permits orderly stopping / shutdown of a group of goroutines via methods which mimic stop of a sync.WaitGroup. There
is a default Terminator accessible through top level functions (Add, Done, Wait and so on) that call the corresponding
Terminator methods.

Another group of functions support logging and tracing. Debug, Error, Info and Warn operate like their package slog equivalents, with
the level of logging modifiable using SetLevel. A custom logging level (LevelTrace) can be supplied to SetLevel to enable tracing. Tracing can be
unconditional when calling Trace, or only enabled for pre-defined identifiers when calling TraceID. Identifiers for TraceID are registered
by calling SetTraceIDs. By default, all debug, error, info and warn messages go to Stdout, and traces go to Stderr; these destinations can
be changed by calling RedirectNormal and RedirectTrace respectively.

Type Set defines methods for manipulating a generic set data structure via the expected operations Add, Contains, Intersection, Members, String
and Union.
*/
package echidna

//go:generate go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest ./... --output body.md
//go:generate ./make_doc.sh
