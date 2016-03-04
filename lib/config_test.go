// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"encoding/json"
	"os"
	"testing"
)

func TestInitialize(t *testing.T) {
	// Prepare the filesystem.
	errCheck(t, os.RemoveAll("/tmp/td"))
	errCheck(t, os.MkdirAll("/tmp/td", 0755))
	errCheck(t, os.MkdirAll("/tmp/td/mordor", 0000))

	// We cannot settle in Mordor.
	errCheck(t, os.Setenv("TD", "/tmp/td/mordor"))
	Initialize()
	_, err := os.Stat("/tmp/td/mordor/old")
	if err == nil {
		t.Fatalf("One does not simply walk into Mordor")
	}

	// A normal setup looks like this.
	dirName = "td"
	errCheck(t, os.Setenv("TD", "/tmp"))
	Initialize()
	_, err = os.Stat("/tmp/td/config.json")
	if err != nil {
		t.Fatalf("Did not expect to encounter error: %v", err)
	}

	errCheck(t, os.Setenv("TD", ""))
	dirName = ".td"
}

func TestInitFS(t *testing.T) {
	// Prepare the filesystem.
	errCheck(t, os.RemoveAll("/tmp/td"))
	errCheck(t, os.MkdirAll("/tmp/td/inside", 0755))
	errCheck(t, os.MkdirAll("/tmp/td/mordor", 0000))

	// You cannot enter Mordor.
	errCheck(t, os.Setenv("TD", "/tmp/td/mordor"))
	err := initFS()
	if err == nil {
		t.Fatalf("One does not simply walk into Mordor")
	}

	// Let's try to do it in a nicer place.
	errCheck(t, os.Setenv("TD", "/tmp/td/inside"))
	errCheck(t, initFS())

	// There's no problem to do that again, directories that already exist are
	// handled gracefully.
	errCheck(t, initFS())

	errCheck(t, os.Setenv("TD", ""))
	dirName = ".td"
}

func TestInitConfig(t *testing.T) {
	// Prepare the filesystem.
	errCheck(t, os.RemoveAll("/tmp/td"))
	errCheck(t, os.MkdirAll("/tmp/td", 0755))
	errCheck(t, os.MkdirAll("/tmp/td/hodor", 0000))
	errCheck(t, os.Setenv("TD", "/tmp"))

	// First of all, let's try to open the config file that is in Mordor.
	config = &configuration{}
	dirName = "td/hodor"
	initConfig()
	if config.logged {
		t.Fatalf("It shouldn't be logged in!")
	}

	// The config file is inside a subdirectory in Mordor that doesn't exist.
	config = &configuration{}
	dirName = "td/hodor/inside"
	initConfig()
	if config.logged {
		t.Fatalf("It shouldn't be logged in!")
	}

	// configFile creates the file, so it's not a valid one.
	config = &configuration{}
	dirName = "td"
	initConfig()
	if config.logged {
		t.Fatalf("It shouldn't be logged in!")
	}
	_, err := os.Stat("/tmp/td/config.json")
	if err != nil {
		t.Fatalf("Did not expect to encounter error: %v", err)
	}

	// A valid config.
	save := &configuration{
		Server: "server",
		Token:  "token",
	}
	body, _ := json.Marshal(save)
	file, _ := os.Create("/tmp/td/config.json")
	defer func() { _ = file.Close() }()
	_, err = file.Write(body)
	errCheck(t, err)

	// And not it should get the previously stored config.
	config = &configuration{}
	initConfig()
	if config.Server != "server" {
		t.Fatalf("Expected %v; Got %v", config.Server, "server")
	}
	if config.Token != "token" {
		t.Fatalf("Expected %v; Got %v", config.Token, "token")
	}
	if !config.logged {
		t.Fatalf("It should be logged in!")
	}

	errCheck(t, os.Setenv("TD", ""))
	dirName = ".td"
}
