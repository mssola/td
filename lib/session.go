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

type configuration struct {
	Server string `json:"server"`
	Token  string `json:"token"`
}

var (
	config *configuration
)

func configFile() (string, error) {
	// Every single command will reach this point eventually, so it's safe to
	// initialize the configuration here.
	config = &configuration{}

	// Create the config file if it doesn't exist yet.
	cfg := filepath.Join(home(), dirName, fileName)
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		dir := filepath.Dir(cfg)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", newError("config file could not be read")
		}

		// Create the config file.
		file, err := os.Create(cfg)
		defer file.Close()
		if err != nil {
			return "", newError("config file could not be read")
		}
	} else if err != nil {
		return "", newError("config file could not be read")
	}
	return cfg, nil
}

func performLogin() error {
	var r request

	fmt.Print("Server: ")
	fmt.Scanf("%v", &config.Server)
	fmt.Print("Name: ")
	fmt.Scanf("%v", &r.Name)
	fmt.Print("Password: ")
	r.Password = string(gopass.GetPasswdMasked())

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
	filePath, err := configFile()
	if err != nil {
		return false
	}

	contents, e := ioutil.ReadFile(filePath)
	if e != nil {
		return false
	}

	if e = json.Unmarshal(contents, &config); e != nil {
		return false
	}
	return config.Token != ""
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
	return fromError(os.RemoveAll(cfg))
}
