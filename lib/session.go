// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
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
)

// loginRequest contains the parameters being used for logging in a user.
type loginRequest struct {
	Name     string
	Password string
}

// LoggedIn returns whether the current user is logged in or not.
func LoggedIn() bool {
	return config.logged
}

// performLogin performs the HTTP request to log in the given user. This
// function assumes that the configuration has already been updated with the
// server to be used.
func performLogin(username, password string) error {
	r := loginRequest{username, password}

	body, _ := json.Marshal(&r)
	reader := bytes.NewReader(body)
	res, err := safeResponse("POST", "/login", reader, false)
	if err != nil {
		return NewError("could not log user in: " + err.Error())
	}
	if res.StatusCode == http.StatusBadRequest {
		return NewError("could not log user in: wrong credentials")
	}

	all, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(all, &config); err != nil {
		return NewError("could not log user in: " + err.Error())
	}
	if config.Token == "" {
		return NewError("could not log user in: no token was given")
	}
	return nil
}

// Login performs the login command.
func Login(server, username, password string) error {
	// Perform the login itself.
	config.Server = server
	if err := performLogin(username, password); err != nil {
		return err
	}

	// Save the configuration.
	fmt.Printf("\nLogging in... ")
	if err := saveConfig(); err != nil {
		return err
	}

	// And fetch the topics.
	fmt.Printf("OK\nFetching topics...\n")
	return fetch()
}

// Logout removes the `.td` directory and everything inside of it.
func Logout() error {
	cfg := filepath.Join(home(), dirName)
	_ = os.RemoveAll(cfg)
	config.logged = false
	return nil
}
