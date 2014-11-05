// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper used to test custom errors. This is used throughout the test suite.
func testError(t *testing.T, err error, msg, see string) {
	assert.NotNil(t, err)
	e, ce := err.(*Error)
	assert.True(t, ce)
	assert.Equal(t, e.message, msg)
	assert.Equal(t, e.see, see)
}

func TestError(t *testing.T) {
	assert.Nil(t, fromError(nil))
	assert.Equal(t, fromError(errors.New("a")).Error(),
		"\x1b[1;49;31merror:\x1b[0;m a.\n")
	msg := See("hello", "help").String()
	assert.Equal(t, msg,
		"\x1b[1;49;31merror:\x1b[0;m hello. \x1b[1;49;39mSee:\x1b[0;m 'td help'.\n")
}
