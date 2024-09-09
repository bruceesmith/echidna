// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the Apache
// License that can be found in the LICENSE file.

package echidna

import "errors"

var (
	// Program is the program name
	Program string
	// Version is the program version
	Version string
)

var (
	// ErrConfigOK is returned when --checkcfg finds the configuration in good order
	ErrConfigOK = errors.New("configuration is OK")
	// ErrVersion is returned when --version is given
	ErrVersion = errors.New("version requested")
)
