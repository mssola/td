// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper type that hides the standard output of the current test. Initialize
// it with the testHelper function, and then call teardown at the end of the
// function (defer).
type helper struct {
	oldStdout *os.File
	writer    *os.File
}

func testHelper() helper {
	var h helper
	var r io.Reader

	h.oldStdout = os.Stdout
	r, h.writer, _ = os.Pipe()
	os.Stdout = h.writer

	go func() {
		io.Copy(ioutil.Discard, r)
	}()
	return h
}

func (h helper) teardown() {
	h.writer.Close()
	os.Stdout = h.oldStdout
}

// This type implements the ttyReader interface as defined in the session.go
// file. This way we don't have to deal with stdin and we can setup multiple
// values easily.
type ttyTest struct {
	name, password string
	server         string
}

// Initialize the proper values and return the initialized request.
func (t ttyTest) login() *request {
	var r request

	config.Server = t.server
	r.Name = t.name
	r.Password = t.password
	return &r
}

// The server used for proper login & logout.
func sessionServer(uName, uPass string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.String() == "/login" {
			// The login itself. Just make sure that the name and password
			// parameters are there.
			body, _ := ioutil.ReadAll(r.Body)
			var r request
			json.Unmarshal(body, &r)
			if r.Name == uName && r.Password == uPass {
				fmt.Fprintln(w, "{\"token\":\"1234\"}")
			} else {
				fmt.Fprintln(w, "{}")
			}
		} else {
			// After logging in it will try to fetch the topics.
			t := []Topic{Topic{Id: "1", Name: "topic"}}
			b, _ := json.Marshal(&t)
			fmt.Fprintln(w, string(b))
		}
	}))
}

func TestAlreadyLoggedIn(t *testing.T) {
	config = &configuration{logged: true}
	testError(t, Login(), "you are already logged in", "logout")
}

func TestServerDown(t *testing.T) {
	h := testHelper()
	defer h.teardown()

	config = &configuration{}
	tty := ttyTest{
		server:   "",
		name:     "name",
		password: "1234",
	}

	err := handleLogin(tty)
	testError(t, err, "could not log user in", "")
}

func TestServerBadJson(t *testing.T) {
	h := testHelper()
	defer h.teardown()

	config = &configuration{}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	tty := ttyTest{
		server:   ts.URL,
		name:     "name",
		password: "1234",
	}

	err := handleLogin(tty)
	testError(t, err, "could not log user in", "")
}

func TestServerLogin(t *testing.T) {
	os.RemoveAll("/tmp/td")
	os.MkdirAll("/tmp/td", 0755)
	os.Setenv("TD", "/tmp")
	dirName = "td"

	h := testHelper()
	defer h.teardown()

	uName, uPass := "name", "1234"

	config = &configuration{}

	// Our test server.
	ts := sessionServer(uName, uPass)
	defer ts.Close()

	// Bad name.
	tty1 := ttyTest{
		server:   ts.URL,
		name:     uName + "o",
		password: uPass,
	}

	err := handleLogin(tty1)
	testError(t, err, "could not log user in", "")

	// Right name.
	tty2 := ttyTest{
		server:   ts.URL,
		name:     uName,
		password: uPass,
	}

	err = handleLogin(tty2)
	assert.Nil(t, err)

	// Check the file system.
	var c configuration
	b, _ := ioutil.ReadFile("/tmp/td/config.json")
	json.Unmarshal(b, &c)
	assert.Equal(t, c.Server, ts.URL)
	assert.Equal(t, c.Token, "1234")

	// Tearing down.
	os.Setenv("TD", "")
	dirName = ".td"
}

func TestLogout(t *testing.T) {
	os.RemoveAll("/tmp/td1")
	os.MkdirAll("/tmp/td1", 0755)
	os.Setenv("TD", "/tmp")
	dirName = "td1"

	h := testHelper()
	defer h.teardown()

	config = &configuration{}

	ts := sessionServer("name", "1234")
	defer ts.Close()

	// First we login.
	tty1 := ttyTest{
		server:   ts.URL,
		name:     "name",
		password: "1234",
	}
	err := handleLogin(tty1)
	assert.Nil(t, err)

	// Test that it's really logged in.
	Initialize()
	assert.True(t, config.logged)

	// And now we logout.
	Logout()

	// Check the FS.
	_, err = os.Stat("/tmp/td1")
	assert.NotNil(t, err)
	assert.True(t, os.IsNotExist(err))

	// If we initialize again, it will say that we're not logged in.
	Initialize()
	assert.False(t, config.logged)

	// Tearing down.
	os.Setenv("TD", "")
	dirName = ".td"
}
