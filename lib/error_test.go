// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

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
