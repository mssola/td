// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mssola/dym"
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
func usage() {
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
	os.Exit(0)
}

// The user passed a wrong number of arguments.
func wrongArguments() {
	var n, msg string

	switch os.Args[1] {
	case "create":
		fallthrough
	case "delete":
		n = "1 argument"
	case "rename":
		n = "2 arguments"
	}
	msg = fmt.Sprintf("the '%v' command requires %v, %v given",
		os.Args[1], n, len(os.Args)-2)
	cmd(lib.See(msg, "--help"))
}

// Show a rather verbose help message, with command suggestions.
func verboseHelp(logged bool) {
	// No command given, just get out of here.
	if len(os.Args) == 1 {
		cmd(lib.See("you are not logged in", "login"))
	}

	e := fmt.Sprintf("'%v' is not a td command", os.Args[1])
	msg := lib.See(e, "--help")

	// Get the commands that are close to the one given by the user.
	d := []string{"login", "--help", "--version", "fetch", "list", "logout",
		"push", "create", "delete", "rename"}
	similars := dym.Similar(d, os.Args[1])

	if len(similars) == 0 {
		fmt.Printf("%v", msg)
	} else {
		str := fmt.Sprintf("%v\nDid you mean one of these?\n", msg)
		for _, v := range similars {
			// Check for an exact match for the given parameter. If there is a
			// match, then it means that the error to be shown is that the user
			// wasn't logged in.
			if v == os.Args[1] {
				if logged {
					wrongArguments()
				} else {
					cmd(lib.See("you are not logged in", "login"))
				}
			}

			str += fmt.Sprintf("\t%v\n", v)
		}
		fmt.Print(str)
	}

	// If we reach this point, then there was no exact match, so it's safe to
	// just quit the program.
	os.Exit(1)
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
	lib.Initialize()

	// All the actions that a non-logged in user can perform.
	if largs == 2 {
		switch os.Args[1] {
		case "login":
			cmd(lib.Login())
		case "--help":
			usage()
		case "--version":
			version()
		}
	}

	// All the other commands require the user to be logged in.
	if !lib.LoggedIn() {
		verboseHelp(false)
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
			verboseHelp(true)
		}
	} else if largs == 3 {
		if os.Args[1] == "create" {
			cmd(lib.Create(os.Args[2]))
		} else if os.Args[1] == "delete" {
			cmd(lib.Delete(os.Args[2]))
		} else {
			verboseHelp(true)
		}
	} else if largs == 4 && os.Args[1] == "rename" {
		cmd(lib.Rename(os.Args[2], os.Args[3]))
	} else {
		verboseHelp(true)
	}
}
