#!/bin/sh
# Copyright (C) 2016 Miquel Sabaté Solà <mikisabate@gmail.com>
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.

# I know it's a bit too much to create an entrypoint.sh file just for setting
# this thing as it should be (I use Ctrl+s to save on Vim, yeah, I'm crazy),
# but I didn't know any better way to do it...
stty -ixon

td $@
