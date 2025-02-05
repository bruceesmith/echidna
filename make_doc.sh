#!/bin/bash

# Copyright © 2024-2025 Bruce Smith <bruceesmith@gmail.com>
# Use of this source code is governed by the MIT
# License that can be found in the LICENSE file.

cat header.md body.md >README.md
rm body.md
sed -i '/^# echidna/d' README.md
sed -i 's/^# program/# <a name="programme">4. program<\/a>/' README.md
sed -i 's/^# stack/# <a name="stack">6. stack<\/a>/' README.md
sed -i 's/^# terminator/# <a name="terminator">7. terminator<\/a>/' README.md
