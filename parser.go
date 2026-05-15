// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package echidna

/*

These functions are based upon functions from https://github.com/urfave/sflags.git and
so reproduce the licensing of that package below

==================================================================================

BSD 3-Clause License

Copyright (c) 2016, Slava Bakhmutov
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

* Neither the name of the copyright holder nor the names of its
  contributors may be used to endorse or promote products derived from
  this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"reflect"
	"strings"
)

// camelToFlag transform s from CamelCase to flag-case
func camelToFlag(s, flagDivider string) string {
	splitted := split(s)
	return strings.ToLower(strings.Join(splitted, flagDivider))
}

// parseField returns the name of a command-line flag for the field argument.
// The makeup of the name matches that used by urfave/sfags/Parse*** functions
func parseField(field reflect.StructField, prefix string) string {
	ignorePrefix := false
	name := camelToFlag(field.Name, opts.divider)
	if tags := strings.Split(field.Tag.Get(opts.tag), ","); len(tags) > 0 {
		switch fName := tags[0]; fName {
		case "-":
			return ""
		case "":
		default:
			fNameSplitted := strings.Split(fName, " ")
			if len(fNameSplitted) > 1 {
				fName = fNameSplitted[0]
			}
			if strings.HasPrefix(fName, "~") {
				name = fName[1:]
				ignorePrefix = true
			} else {
				name = fName
			}
		}

	}

	if prefix != "" && !ignorePrefix {
		name = prefix + name
	}
	return name
}
