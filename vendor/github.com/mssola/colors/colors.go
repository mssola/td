// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package colors

import (
	"fmt"
)

// Represents a displayable color.
type Colors int

const (
	Black Colors = iota
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White

	// "Saved" preserves the color being used by the color schema.
	Saved = 9
)

// The mode in which the message should be displayed.
type Mode int

const (
	Regular Mode = iota
	Bold
	Underlined = 4
	Reversed   = 7
)

// This is the struct to be built in order to colorize strings.
type Color struct {
	// The foreground color.
	Foreground Colors

	// The background color.
	Background Colors

	// The mode of display.
	Mode Mode
}

// Returns a new Color object with the default values. The default values are
// "Saved" for both colors and "Regular" for the mode.
func Default() *Color {
	return &Color{
		Foreground: Saved,
		Background: Saved,
		Mode:       Regular,
	}
}

// Change the color. If bg is true, then the given color will be applied to the
// background. Otherwise the foreground will be the one modified.
func (c *Color) Change(color Colors, bg bool) {
	if bg {
		c.Background = color
	} else {
		c.Foreground = color
	}
}

// Set the given mode to this Color object.
func (c *Color) SetMode(m Mode) {
	c.Mode = m
}

// Colorizes the given string with the current settings.
func (c *Color) Get(s string) string {
	return fmt.Sprintf("\x1b[%v;4%v;3%vm%v\x1b[0;m", c.Mode, c.Background,
		c.Foreground, s)
}
