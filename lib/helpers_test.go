// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"log"
	"os"
	"path/filepath"
	"testing"
)

// getWd returns the current working directory where the test environment
// should reside.
func getWd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not get working directory: %v", err)
	}
	if filepath.Base(wd) == "lib" {
		wd = filepath.Dir(wd)
		if err := os.Chdir(wd); err != nil {
			log.Fatalf("Could not cd into working directory: %v", err)
		}
	}
	return wd
}

func startTestEnv(t *testing.T) {
	config = &configuration{}
	Insecure = true

	wd := getWd()
	path := filepath.Join(wd, dirName)

	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("Could not cleanup test environment: %v", err)
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		t.Fatalf("Could not create test environment: %v", err)
	}
	if err := os.Setenv("TD", wd); err != nil {
		_ = os.RemoveAll(dirName)
		t.Fatalf("Could not create test environment: %v", err)
	}
}

func stopTestEnv(t *testing.T) {
	path := filepath.Join(getWd(), dirName)

	if err := os.RemoveAll(path); err != nil {
		t.Fatalf("Could not cleanup test environment: %v", err)
	}
	if err := os.Setenv("TD", ""); err != nil {
		t.Fatalf("Could not create test environment: %v", err)
	}
}

func errCheck(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Did not expect to get an error: %v", err)
	}
}
