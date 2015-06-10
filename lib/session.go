// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"code.google.com/p/go.crypto/ssh/terminal"
)

// The struct representing the login request.
type request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// This is an interface that just defines one method: login. Implementers of
// this interface handle how the login data gets collected. This is useful
// because in the application we collect this data through the TTY, but for
// the tests we set them directly.
type loginReader interface {
	// This method is responsible for giving a value to the Server field of the
	// configuration struct, and it also initializes the login request, that is
	// the value being returned by this function.
	login() *request
}

// We don't need any extra fields, this type only wants to implement the
// loginReader interface.
type tty struct{}

// This application implements the login function by reading the different
// values from the TTY.
func (tty) login() *request {
	var r request

	fmt.Print("Server: ")
	fmt.Scanf("%v", &config.Server)
	fmt.Print("Name: ")
	fmt.Scanf("%v", &r.Name)
	fmt.Print("Password: ")
	b, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
	r.Password = string(b)
	return &r
}

// Read the login data from the given loginReader and perform the login command.
func handleLogin(t loginReader) error {
	if LoggedIn() {
		return See("you are already logged in", "logout")
	}

	// Get the initial values for the request.
	r := t.login()
	fmt.Printf("\n")

	// And perform the login itself.
	if err := performLogin(r); err != nil {
		fmt.Printf("\nLogging in... FAIL\n")
		return err
	}

	// Save config
	fmt.Printf("\nLogging in... ")
	body, _ := json.Marshal(config)
	filePath, _ := configFile()
	f, _ := os.Create(filePath)
	f.Write(body)
	f.Close()
	fmt.Printf("OK\nFetching topics...\n")
	return Fetch()
}

// Called by the handleLogin function, it performs the HTTP request to log in
// the given user. This function also deals with response, by initializing the
// config global variable accordingly.
func performLogin(r *request) error {
	url := requestUrl("/login", false)
	body, _ := json.Marshal(r)
	reader := bytes.NewReader(body)
	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		return newError("could not log user in")
	}

	all, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(all, &config); err != nil {
		return newError("could not log user in")
	}
	if config.Token == "" {
		return newError("could not log user in")
	}
	return nil
}

// Returns true if it has been detected that the current user is logged in.
func LoggedIn() bool {
	return config.logged
}

// The login command.
func Login() error {
	// We handle the login with the values provided through the TTY.
	return handleLogin(tty{})
}

// Logout: remove the `.td` directory and everything inside of it.
func Logout() error {
	cfg := filepath.Join(home(), dirName)
	os.RemoveAll(cfg)
	return nil
}
