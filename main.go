// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mssola/td/lib"
)

const (
	// The major element for the version number.
	major = 0

	// The minor element for the version number.
	minor = 1

	// The patch level for the version number.
	patch = 0
)

// Show the usage string. The parameter will be used as the exit status code.
func usage(status int) {
	msg := []string{
		"usage: td [--version | --help | <command>] [args]",
		"",
		"The available commands are:",
		"  create\tCreate a new topic. It expects one extra argument: the name.",
		"  delete\tDelete a topic. It expects one extra argument: the name.",
		"  fetch \tFetch all the info from the server.",
		"  list  \tList the available topics.",
		"  logout\tDelete the current session.",
		"  push  \tPush all the local info to the server.",
		"  rename\tRename a topic. You have to pass the old name and " +
			"the new one.",
	}
	fmt.Printf("%v\n", strings.Join(msg, "\n"))
	os.Exit(status)
}

// Show the version of this program.
func version() {
	if patch == 0 {
		fmt.Printf("td version %v.%v\n", major, minor)
	} else {
		fmt.Printf("td version %v.%v.%v\n", major, minor, patch)
	}
	os.Exit(0)
}

// All the commands return an error value. This function evaluates this error
// and prints an error message if it's required.
func cmd(err error) {
	if err == nil {
		os.Exit(0)
	}
	fmt.Print(err)
	os.Exit(1)
}

func main() {
	largs := len(os.Args)

	// All the actions that a non-logged in user can perform.
	if largs == 2 {
		switch os.Args[1] {
		case "login":
			cmd(lib.Login())
		case "--help":
			usage(0)
		case "--version":
			version()
		}
	}

	// All the other commands require the user to be logged in.
	if !lib.LoggedIn() {
		cmd(errors.New("you are not logged in.\nTry: `td login`"))
	}

	// Let's execute the given command now.
	if largs == 1 {
		cmd(lib.Edit())
	} else if largs == 2 {
		switch os.Args[1] {
		case "fetch":
			cmd(lib.Fetch())
		case "list":
			cmd(lib.List())
		case "logout":
			cmd(lib.Logout())
		case "push":
			cmd(lib.Push())
		default:
			usage(1)
		}
	} else if largs == 3 {
		if os.Args[1] == "create" {
			cmd(lib.Create(os.Args[2]))
		} else if os.Args[1] == "delete" {
			cmd(lib.Delete(os.Args[2]))
		} else {
			usage(1)
		}
	} else if largs == 4 && os.Args[1] == "rename" {
		cmd(lib.Rename(os.Args[2], os.Args[3]))
	} else {
		usage(1)
	}
}
