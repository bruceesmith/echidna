// Copyright © 2024 Bruce Smith <bruceesmith@gmail.com>
// Use of this source code is governed by the MIT
// License that can be found in the LICENSE file.

/*
Package echidna provides sub-packages for building robust Go daemons and CLIs

  - program builds upon the Github packages [knadh/koanf] and [urfave/cli/v3] to make it extremely simple to use the
    features of those two excellent packages in concert.

Refer to the documentation for the individual packages for more details.

[urfave/cli/v3]: https://github.com/urfave/cli
[knadh/koanf]: https://github.com/knadh/koanf
*/
package echidna

//go:generate go run github.com/princjef/gomarkdoc/cmd/gomarkdoc@latest ./... --output body.md
//go:generate ./make_doc.sh
