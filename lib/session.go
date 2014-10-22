// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
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

	"github.com/howeyc/gopass"
)

type request struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

func performLogin() error {
	var r request

	fmt.Print("Server: ")
	fmt.Scanf("%v", &config.Server)
	fmt.Print("Name: ")
	fmt.Scanf("%v", &r.Name)
	fmt.Print("Password: ")
	r.Password = string(gopass.GetPasswd())

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

func LoggedIn() bool {
	return config.logged
}

// TODO: it doesn't actually fetch after login :/
func Login() error {
	if LoggedIn() {
		return See("you are already logged in", "logout")
	}

	// Get the initial values for the request.
	if err := performLogin(); err != nil {
		fmt.Printf("\nLogging in... FAIL\n")
		return err
	}

	// Save config
	fmt.Printf("\nLogging in... ")
	body, _ := json.Marshal(config)
	filePath, _ := configFile()
	f, _ := os.Create(filePath)
	defer f.Close()
	if _, err := f.Write(body); err != nil {
		fmt.Printf("FAIL\n")
		return newError("could not save the config!")
	}
	fmt.Printf("OK\n")
	fmt.Printf("Fetching topics...\n")
	return Fetch()
}

// TODO: somewhere in this code is printing <nil> on success :/
func Logout() error {
	// Remove the `.td` directory and everything inside of it.
	cfg := filepath.Join(home(), dirName)
	if err := os.RemoveAll(cfg); err != nil {
		return fromError(err)
	}
	return nil
}
