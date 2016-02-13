// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/mssola/capture"
)

// The server used for proper login & logout.
func sessionServer(username, password string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/login" {
			// The login itself. Just make sure that the name and password
			// parameters are there.
			body, _ := ioutil.ReadAll(r.Body)
			var r loginRequest
			_ = json.Unmarshal(body, &r)
			if r.Name == username && r.Password == password {
				fmt.Fprintln(w, "{\"token\":\"1234\"}")
			} else if r.Password == username {
				// Let's simulate this kind of malformed error.
				fmt.Fprintln(w, "{\"token\":\"\"}")
			} else {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "{}")
			}
		} else {
			// After logging in it will try to fetch the topics.
			t := []Topic{{ID: "1", Name: "topic"}}
			b, _ := json.Marshal(&t)
			fmt.Fprintln(w, string(b))
		}
	}))
}

func TestServerDown(t *testing.T) {
	config = &configuration{}
	if err := Login("", "name", "1234"); err == nil {
		t.Fatalf("Expected an error to occur")
	}
}

func TestServerBadJson(t *testing.T) {
	Insecure = true
	config = &configuration{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	err := Login(ts.URL, "name", "1234")
	if err == nil {
		t.Fatalf("It should have failed!")
	}
	msg := "could not log user in: invalid character 'H' looking for beginning of value"
	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("Got '%v'; expected %s", err, msg)
	}
}

func testLogin(t *testing.T, url, username, password string) {
	var err error

	capture.All(func() { err = Login(url, username, password) })
	if err != nil {
		t.Fatalf("Should not given an error: %v", err)
	}

	// Check the file system.
	var c configuration
	wd := os.Getenv("TD")
	b, _ := ioutil.ReadFile(filepath.Join(wd, dirName, configName))
	errCheck(t, json.Unmarshal(b, &c))
	if c.Server != url {
		t.Fatalf("Got: %v; Expected: %v", c.Server, url)
	}
	if c.Token != "1234" {
		t.Fatalf("Got: '%v'; Expected: '1234'", c.Token)
	}
}

func TestServerLogin(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	username, password := "name", "1234"
	ts := sessionServer(username, password)
	defer ts.Close()

	// Bad name.
	err := Login(ts.URL, username+"o", password)
	if err == nil {
		t.Fatalf("We were expecting an error here...")
	}
	msg := "could not log user in: wrong credentials"
	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("Got: %v; expected: %v", err.Error(), msg)
	}

	// Use the common helper.
	testLogin(t, ts.URL, username, password)
}

func TestServerBadLogin(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	username := "name"
	ts := sessionServer(username, "1234")
	defer ts.Close()

	// Bad password, this mock will return an empty token.
	msg := "could not log user in: no token was given"
	if err := Login(ts.URL, username, username); err == nil {
		t.Fatalf("Should have given an error!")
	} else if !strings.Contains(err.Error(), msg) {
		t.Fatalf("Got: %v; expected: %v", err, msg)
	}
}

func TestCannotLoginFS(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	if err := os.Setenv("TD", "/llalala"); err != nil {
		t.Fatalf("Could not create test environment: %v", err)
	}
	defer func() { _ = os.Setenv("TD", "") }()

	username, password := "name", "1234"
	ts := sessionServer(username, password)
	defer ts.Close()

	msg := "no such file or directory"
	var err error

	capture.All(func() { err = Login(ts.URL, username, password) })
	if err == nil {
		t.Fatalf("Should have given an error!")
	} else if !strings.Contains(err.Error(), msg) {
		t.Fatalf("Got: %v; expected: %v", err, msg)
	}
}

func TestLogout(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	username, password := "name", "1234"
	ts := sessionServer(username, password)
	defer ts.Close()

	// Login & logout
	testLogin(t, ts.URL, username, password)
	if err := Logout(); err != nil {
		t.Fatalf("Should not given an error: %v", err)
	}
	if LoggedIn() {
		t.Fatalf("It says that it's logged in when it's not!")
	}
	cfg := filepath.Join(home(), dirName)
	if _, err := os.Stat(cfg); err == nil || !os.IsNotExist(err) {
		t.Fatalf("Oops: %v", err)
	}
}
