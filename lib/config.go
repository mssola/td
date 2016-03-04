// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

type configuration struct {
	Server string `json:"server"`
	Token  string `json:"token"`
	logged bool
}

const (
	configName = "config.json"
)

var (
	config *configuration
)

// Initialize performs the needed initialization for the application.
func Initialize() {
	config = &configuration{logged: false}

	// Check out the file system. We do this so we can make sure that
	// any following command touching the file system can do it safely.
	if err := initFS(); err == nil {
		// And initialize the "config" global variable.
		initConfig()
	}
}

func initFS() error {
	if err := checkDir(oldDir); err != nil {
		return err
	}
	if err := checkDir(newDir); err != nil {
		return err
	}
	return checkDir(tmpDir)
}

func checkDir(dir string) error {
	s := filepath.Join(home(), dirName, dir)
	if _, err := os.Stat(s); err != nil {
		if os.IsNotExist(err) {
			_ = os.MkdirAll(s, 0755)
			return nil
		}
		return NewError("there's something wrong with the file system")
	}

	// The directory exists, try to give it some cool permissions.
	return os.Chmod(s, 0755)
}

func initConfig() {
	// Try to get the config file. If that's not possible, then it means that
	// the user is not logged in.
	filePath, err := configFile()
	if err != nil {
		config.logged = false
		return
	}

	// And finally we'll initialize the config variable properly.
	contents, _ := ioutil.ReadFile(filePath)
	if err = json.Unmarshal(contents, &config); err != nil {
		config.logged = false
		return
	}
	config.logged = (config.Token != "")
	return
}

func configFile() (string, error) {
	// Create the config file if it doesn't exist yet.
	cfg := filepath.Join(home(), dirName, configName)
	if _, err := os.Stat(cfg); os.IsNotExist(err) {
		dir := filepath.Dir(cfg)
		_ = os.MkdirAll(dir, 0755)

		// Create the config file.
		file, _ := os.Create(cfg)
		_ = file.Close()
	} else if err != nil {
		return "", NewError("config file could not be read")
	}
	return cfg, nil
}

func saveConfig() error {
	body, _ := json.Marshal(config)
	filePath, _ := configFile()

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	_, _ = f.Write(body)
	_ = f.Close()
	return nil
}
