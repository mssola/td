
# Colors [![Build Status](https://travis-ci.org/mssola/colors.png?branch=master)](https://travis-ci.org/mssola/colors) [![GoDoc](https://godoc.org/github.com/mssola/colors?status.png)](http://godoc.org/github.com/mssola/colors)

## About this

This package offers a simple way to colorize strings that will be printed into
the standard output. Note that Windows is not supported. The available colors
are:

    Black
    Red
    Green
    Yellow
    Blue
    Magenta
    Cyan
    White
    Saved

The `Saved` color asks to preserve the color being used in the console.
Moreover, this packages acknowledges different modes:

    Regular
    Bold
    Underlined
    Reversed

In order to use all this, we'll have to create a new `Color` object. You
usually want to do this with the `Default` function. This function will build a
new `Color` object with the default values. With a `Color` object we can call
the following methods:

```go
// Change the color. The bg parameter asks whether this color has to be
// applied to the background or not.
func (c *Color) Change(color Colors, bg bool)

// Changes the mode for the given one.
func (c *Color) SetMode(m Mode)

// Returns the given string but with the color & mode applied.
func (c *Color) Get(s string) string
```

For more information check out the [documentation](http://godoc.org/github.com/mssola/colors).

## tl;dr

```go
package main

import (
    "fmt"

    "github.com/mssola/colors"
)

func main() {
    // By default it respects the current terminal schema.
    c := colors.Default()
    fmt.Printf(c.Get("Nothing\n"))

    // c := colors.Default() is equivalent to:
    //
    // c := &colors.Color{
    //    Foreground: colors.Saved,
    //    Background: colors.Saved,
    //    Mode: colors.Regular,
    // }
    //
    // You might want to initialize the struct directly if you just want to
    // print a colorful message and get out.

    // Let's print a red bold message.
    c.Change(colors.Red, false)
    c.SetMode(colors.Bold)
    fmt.Printf(c.Get("Red\n"))

    // It will print the same as before but reversed: the background is now red
    // and the foreground white.
    c.SetMode(colors.Reversed)
    fmt.Printf(c.Get("Reversed Red\n"))
}
```

## License

Copyright &copy; 2014-2016 Miquel Sabaté Solà, released under the MIT License.

