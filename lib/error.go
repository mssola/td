// Copyright (C) 2014-2017 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"fmt"

	"github.com/mssola/colors"
)

// Error type to be used throughout this application.
type Error struct {
	message string
	see     string
}

// NewError builds a new error from the given message.
func NewError(message string) *Error {
	return &Error{message: message, see: ""}
}

// See builds a new error from the given messages.
func See(message, see string) *Error {
	return &Error{message: message, see: see}
}

// fromError build a new custom error from the given standard error.
func fromError(err error) *Error {
	if err == nil {
		return nil
	}
	return NewError(err.Error())
}

// So we implement the Stringer interface.
func (e *Error) String() string {
	red := &colors.Color{
		Foreground: colors.Red,
		Background: colors.Saved,
		Mode:       colors.Bold,
	}
	white := colors.Default()
	white.SetMode(colors.Bold)

	str := fmt.Sprintf("%v %v.", red.Get("error:"), e.message)
	if e.see != "" {
		str += fmt.Sprintf(" %v 'td %v'.", white.Get("See:"), e.see)
	}
	return str + "\n"
}

// So we implement the error interface.
func (e *Error) Error() string {
	return e.String()
}

// Print a warning to stdout.
func warning(str, extra string) {
	red := &colors.Color{
		Foreground: colors.Red,
		Background: colors.Saved,
		Mode:       colors.Bold,
	}
	str = "%v: " + str + "\n"
	if extra == "" {
		fmt.Printf(str, red.Get("warning"))
	} else {
		fmt.Printf(str, red.Get("warning"), extra)
	}
}
