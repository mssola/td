// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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

const (
	dirName  = ".td"
	fileName = ".config.json"
)

var (
	config *configuration
)

func configFile() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		panic("You don't have the $HOME environment variable set!")
	}

	// Every single command will reach this point eventually, so it's safe to
	// initialize the configuration here.
	config = &configuration{}

	// Create the config file if it doesn't exist yet.
	cfg := filepath.Join(home, dirName, fileName)
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		dir := filepath.Dir(cfg)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return "", errors.New("config file could not be read!")
		}

		// Create the config file.
		file, err := os.Create(cfg)
		defer file.Close()
		if err != nil {
			return "", errors.New("config file could not be read!")
		}
	} else if err != nil {
		return "", errors.New("config file could not be read!")
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

	url := config.Server
	if !strings.HasSuffix(url, "/") {
		url += "/"
	}
	url += "login"

	body, _ := json.Marshal(r)
	reader := bytes.NewReader(body)
	res, err := http.Post(url, "application/json", reader)
	if err != nil {
		return errors.New("Could not log user in!")
	}

	all, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(all, &config); err != nil {
		return errors.New("Could not log user in!")
	}
	if config.Token == "" {
		return errors.New("Could not log user in!")
	}
	return nil
}

func LoggedIn() bool {
	filePath, err := configFile()
	if err != nil {
		return false
	}

	contents, err := ioutil.ReadFile(filePath)
	if err != nil {
		return false
	}

	if err := json.Unmarshal(contents, &config); err != nil {
		return false
	}
	return config.Token != ""
}

func Login() error {
	if LoggedIn() {
		return errors.New("you are already logged in.\nTry: `td logout`")
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
		return errors.New("could not save the config!")
	}
	fmt.Printf("OK\n")
	return nil
}

func Logout() error {
	home := os.Getenv("HOME")
	if home == "" {
		panic("You don't have the $HOME environment variable set!")
	}

	// Remove the `.td` directory and everything inside of it.
	cfg := filepath.Join(home, dirName)
	return os.RemoveAll(cfg)
}
