// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

package logger

import (
	"fmt"
	"strings"
)

// Traces is the list of trace IDs enabled
type Traces []string

// String is a convenience method for pflag.Value
func (t *Traces) String() (s string) {
	return fmt.Sprint(*t)
}

// Set is a convenience method for pflag.Value
func (t *Traces) Set(ts string) (err error) {
	parts := strings.Split(ts, ",")
	if len(*t) == 0 {
		*t = make([]string, 0)
	}
	*t = append(*t, parts...)
	return
}

// Type is a conveniene method for pflag.Value
func (t *Traces) Type() string {
	return "Traces"
}
