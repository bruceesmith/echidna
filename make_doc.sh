#!/bin/bash

# Copyright © 2024-2025 Bruce Smith <bruceesmith@gmail.com>
# Use of this source code is governed by the MIT
# License that can be found in the LICENSE file.

cat header.md body.md >README.md
rm body.md
sed -i '/^# echidna/d' README.md
sed -i 's/^# custom/# <a name="custom">2. custom<\/a>/' README.md
sed -i 's/^# logger/# <a name="logger">3. logger<\/a>/' README.md
sed -i 's/^# set/# <a name="set">4. set<\/a>/' README.md
sed -i 's/^# stack/# <a name="stack">4. stack<\/a>/' README.md
sed -i 's/^# terminator/# <a name="terminator">5. terminator<\/a>/' README.md
sed -i 's/^# version/# <a name="version">6. version<\/a>/' README.md
sed -i 's/^# vpr/# <a name="vpr">7. vpr<\/a>/' README.md
