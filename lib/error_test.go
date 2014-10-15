// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	assert.Nil(t, fromError(nil))
	assert.Equal(t, fromError(errors.New("a")).Error(), "td: a.\n")
	msg := See("hello", "help").String()
	assert.Equal(t, msg, "td: hello. See: 'td help'.\n")
}
