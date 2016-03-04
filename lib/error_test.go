// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"errors"
	"testing"
)

func TestError(t *testing.T) {
	if fromError(nil) != nil {
		t.Fatalf("Did not expect to encounter error: %v", fromError(nil))
	}

	msg := "\x1b[1;49;31merror:\x1b[0;m a.\n"
	err := fromError(errors.New("a")).Error()
	if err != msg {
		t.Fatalf("Expected: %v; got: %v", msg, err)
	}

	msg = See("hello", "help").String()
	expected := "\x1b[1;49;31merror:\x1b[0;m hello. \x1b[1;49;39mSee:\x1b[0;m 'td help'.\n"
	if msg != expected {
		t.Fatalf("Expected: %v; got: %v", msg, expected)
	}
}
