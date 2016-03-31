// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/codegangsta/cli"
	"github.com/mssola/td/lib"
)

func version() string {
	const (
		major = 0
		minor = 1
		patch = 0
	)
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}

// errAndExit prints the given error if it was not nil and exits with the
// proper exit code.
func errAndExit(err error) {
	if err == nil {
		os.Exit(0)
	}
	fmt.Print(err)
	os.Exit(1)
}

// flagOrPrompt tries to fetch the value for the requested flag from the given
// CLI context. If that is not possible, then it prompts the user asking for
// the information. If the `secure` parameter is set to true, then the password
// won't be shown.
func flagOrPrompt(ctx *cli.Context, name string, secure bool) (string, error) {
	// If the flag already provides the value, just return it.
	if val := ctx.String(name); val != "" {
		return val, nil
	}

	// Show the prompt by uppercasing the first letter of the given name.
	r, n := utf8.DecodeRuneInString(name)
	fmt.Print(string(unicode.ToUpper(r)) + name[n:] + ": ")

	// And now get the value from the terminal.
	if secure {
		b, err := readPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", err
		}
		fmt.Println()
		return string(b), nil
	}
	var value string
	fmt.Scanf("%s", &value)
	return value, nil
}

// readLoginDetails reads the flags passed to the `login` command in order to
// fetch details needed for the `lib.Login` function. If a flag is not passed,
// then the user will be prompted to give the information manually.
func readLoginDetails(ctx *cli.Context) (string, string, string, error) {
	server, err := flagOrPrompt(ctx, "server", false)
	if err != nil {
		return "", "", "", err
	}
	username, err := flagOrPrompt(ctx, "username", false)
	if err != nil {
		return "", "", "", err
	}
	password, err := flagOrPrompt(ctx, "password", true)
	if err != nil {
		return "", "", "", err
	}
	return server, username, password, nil
}

// loggedCommand wraps the given function by making sure that the current user
// is logged in. If this is not the case, it shows an error message and exits.
func loggedCommand(f func(*cli.Context)) func(*cli.Context) {
	return func(ctx *cli.Context) {
		if !lib.LoggedIn() {
			errAndExit(lib.See("you are not logged in", "login"))
		}
		f(ctx)
	}
}

// require enforces the given cli context to have the expected number of
// arguments.
func require(ctx *cli.Context, expected int) {
	if len(ctx.Args()) != expected {
		msg := fmt.Sprintf("the '%s' command requires %d arguments, %d given",
			ctx.Command.Name, expected, len(ctx.Args()))
		fmt.Println(lib.NewError(msg))
		cli.ShowAppHelp(ctx)
		os.Exit(1)
	}
}

func main() {
	lib.Initialize()

	app := cli.NewApp()
	app.Name = "td"
	app.Usage = "A CLI tool for a 'todo' server."
	app.Version = version()

	app.CommandNotFound = func(context *cli.Context, cmd string) {
		fmt.Printf("Incorrect usage: command '%v' does not exist.\n\n", cmd)
		cli.ShowAppHelp(context)
	}

	app.Action = loggedCommand(func(ctx *cli.Context) {
		errAndExit(lib.Edit())
	})

	app.Commands = []cli.Command{
		{
			Name: "create",
			Usage: "Create a new topic. " +
				"It requires one argument, which is the name of the new topic.",
			Action: loggedCommand(func(ctx *cli.Context) {
				require(ctx, 1)
				errAndExit(lib.Create(ctx.Args()[0]))
			}),
		},
		{
			Name:  "delete",
			Usage: "Delete a topic. It expects one extra argument: the name.",
			Action: loggedCommand(func(ctx *cli.Context) {
				require(ctx, 1)
				errAndExit(lib.Delete(ctx.Args()[0]))
			}),
		},
		{
			Name:   "list",
			Usage:  "List the available topics.",
			Action: loggedCommand(func(ctx *cli.Context) { errAndExit(lib.List()) }),
		},
		{
			Name:  "login",
			Usage: "Log the current user.",
			Action: func(ctx *cli.Context) {
				if lib.LoggedIn() {
					fmt.Println("You are already logged in. Doing nothing...")
					os.Exit(0)
				}

				server, name, password, err := readLoginDetails(ctx)
				if err != nil {
					errAndExit(err)
				}
				if server == "" || name == "" || password == "" {
					errAndExit(lib.NewError("missing information"))
				}
				errAndExit(lib.Login(server, name, password))
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "s, server",
					Usage: "The URL where the 'todo' application is hosted.",
				},
				cli.StringFlag{
					Name:  "u, username",
					Usage: "Username.",
				},
				cli.StringFlag{
					Name:  "p, password",
					Usage: "Password.",
				},
			},
		},
		{
			Name:   "logout",
			Usage:  "Delete the current session.",
			Action: loggedCommand(func(ctx *cli.Context) { errAndExit(lib.Logout()) }),
		},
		{
			Name:  "rename",
			Usage: "Rename a topic. You have to pass the old name and the new one.",
			Action: loggedCommand(func(ctx *cli.Context) {
				require(ctx, 2)
				errAndExit(lib.Rename(ctx.Args()[0], ctx.Args()[1]))
			}),
		},
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:        "insecure",
			Usage:       "Allow the usage of insecure connections.",
			Destination: &lib.Insecure,
		},
		cli.BoolTFlag{
			Name:        "tlsverify",
			Usage:       "Verify the remote server. Ignored if --insecure is set to true.",
			Destination: &lib.TLSVerify,
		},
	}

	app.RunAndExitOnError()
}
