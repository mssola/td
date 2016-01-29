// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInitialize(t *testing.T) {
	// Prepare the filesystem.
	os.RemoveAll("/tmp/td")
	os.MkdirAll("/tmp/td", 0755)
	os.MkdirAll("/tmp/td/mordor", 0000)

	// We cannot settle in Mordor.
	os.Setenv("TD", "/tmp/td/mordor")
	Initialize()
	_, err := os.Stat("/tmp/td/mordor/old")
	assert.NotNil(t, err)

	// A normal setup looks like this.
	dirName = "td"
	os.Setenv("TD", "/tmp")
	Initialize()
	_, err = os.Stat("/tmp/td/config.json")
	assert.Nil(t, err)

	os.Setenv("TD", "")
	dirName = ".td"
}

func TestInitFS(t *testing.T) {
	// Prepare the filesystem.
	os.RemoveAll("/tmp/td")
	os.MkdirAll("/tmp/td/inside", 0755)
	os.MkdirAll("/tmp/td/mordor", 0000)

	// You cannot enter Mordor.
	os.Setenv("TD", "/tmp/td/mordor")
	err := initFS()
	assert.NotNil(t, err)

	// Let's try to do it in a nicer place.
	os.Setenv("TD", "/tmp/td/inside")
	err = initFS()
	assert.Nil(t, err)

	// There's no problem to do that again, directories that already exist are
	// handled gracefully.
	err = initFS()
	assert.Nil(t, err)

	os.Setenv("TD", "")
	dirName = ".td"
}

func TestInitConfig(t *testing.T) {
	// Prepare the filesystem.
	os.RemoveAll("/tmp/td")
	os.MkdirAll("/tmp/td", 0755)
	os.MkdirAll("/tmp/td/hodor", 0000)
	os.Setenv("TD", "/tmp")

	// First of all, let's try to open the config file that is in Mordor.
	config = &configuration{}
	dirName = "td/hodor"
	initConfig()
	assert.False(t, config.logged)

	// The config file is inside a subdirectory in Mordor that doesn't exist.
	config = &configuration{}
	dirName = "td/hodor/inside"
	initConfig()
	assert.False(t, config.logged)

	// configFile creates the file, so it's not a valid one.
	config = &configuration{}
	dirName = "td"
	initConfig()
	assert.False(t, config.logged)
	_, err := os.Stat("/tmp/td/config.json")
	assert.Nil(t, err)

	// A valid config.
	save := &configuration{
		Server: "server",
		Token:  "token",
	}
	body, _ := json.Marshal(save)
	file, _ := os.Create("/tmp/td/config.json")
	defer file.Close()
	file.Write(body)

	// And not it should get the previously stored config.
	config = &configuration{}
	initConfig()
	assert.Equal(t, config.Server, "server")
	assert.Equal(t, config.Token, "token")
	assert.True(t, config.logged)

	os.Setenv("TD", "")
	dirName = ".td"
}
