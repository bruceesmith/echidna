// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

/*
Package echidna provides sub-packages for building robust Go daemons and CLIs

  - observable is an implementation of the [Gang of Four] [observer] pattern, useful in event-driven programs such as GUIs.

  - program builds upon the Github packages [knadh/koanf] and [urfave/cli/v3] to make it extremely simple to use the
    features of those two excellent packages in concert.

  - set defines goroutine-safe methods for manipulating a generic [set] data structure via the standard operations Add, Contains,
    Intersection, Members, String and Union.

  - stack defines goroutine-safe methods for manipulating a generic [stack] data structure via the standard operations IsEmpty,
    Peek, Pop, Push and Size.

  - terminator permits orderly stopping / shutdown of a group of goroutines via methods which mimic a [sync.WaitGroup].
    There is a default [terminator.Terminator] accessible through top level functions (Add, Done, Wait and so on) that call the
    corresponding Terminator methods.

Refer to the documentation for the individual packages for more details.

[urfave/cli/v3]: https://github.com/urfave/cli
[knadh/koanf]: https://github.com/knadh/koanf
[set]: https://en.wikipedia.org/wiki/Set_(abstract_data_type)
[stack]: https://en.wikipedia.org/wiki/Stack_(abstract_data_type)
[Gang of Four]: https://en.wikipedia.org/wiki/Design_Patterns
[observer]: https://en.wikipedia.org/wiki/Observer_pattern
*/
package echidna

//go:generate go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest ./... --output body.md
//go:generate ./make_doc.sh
