// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import "fmt"

// The Error type to be used throughout this application.
type Error struct {
	message string
	see     string
}

// Build a new error from the given message.
func newError(message string) *Error {
	return &Error{
		message: message,
		see:     "",
	}
}

// Build a new error from the given messages.
func See(message, see string) *Error {
	return &Error{
		message: message,
		see:     see,
	}
}

// Build a new custom error from the given standard error.
func fromError(err error) *Error {
	if err == nil {
		return nil
	}
	return newError(err.Error())
}

// So we implement the Stringer interface.
func (e *Error) String() string {
	str := fmt.Sprintf("td: %v.", e.message)
	if e.see != "" {
		str += fmt.Sprintf(" See: 'td %v'.", e.see)
	}
	return str + "\n"
}

// Se we implement the error interface.
func (e *Error) Error() string {
	return e.String()
}
