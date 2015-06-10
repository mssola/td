// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHome(t *testing.T) {
	casa := os.Getenv("HOME")
	assert.Equal(t, casa, home())

	os.Setenv("TD", "/home/lala")
	assert.Equal(t, "/home/lala", home())

	defer func() {
		if err := recover(); err != nil {
			msg := "You don't have the $HOME environment variable set"
			assert.Equal(t, msg, err)
		}
	}()
	os.Setenv("HOME", "")
	os.Setenv("TD", "")
	home()

	os.Setenv("HOME", casa)
}

func TestEditor(t *testing.T) {
	os.Setenv("EDITOR", "emacs")
	assert.Equal(t, editor(), "emacs")
	os.Setenv("EDITOR", "")
	assert.Equal(t, editor(), defaultEditor)
}

func TestCopyFile(t *testing.T) {
	// Prepare the filesystem.
	os.RemoveAll("/tmp/td")
	os.MkdirAll("/tmp/td/good", 0755)
	os.MkdirAll("/tmp/td/mordor", 0000)
	f, err := os.Create("/tmp/td/good/test.txt")
	defer f.Close()
	assert.Nil(t, err)
	f.WriteString("hello world")

	// Good copy of the "good" dir.
	err = copyDir("/tmp/td/good", "/tmp/td/lala")
	assert.Nil(t, err)
	body, err := ioutil.ReadFile("/tmp/td/lala/test.txt")
	assert.Nil(t, err)
	assert.Equal(t, string(body), "hello world")

	// Copying the same directory replaces the old one.
	f.WriteString(" lala")
	err = copyDir("/tmp/td/good", "/tmp/td/lala")
	assert.Nil(t, err)
	body, err = ioutil.ReadFile("/tmp/td/lala/test.txt")
	assert.Nil(t, err)
	assert.Equal(t, string(body), "hello world lala")

	// Trying from a directory that doesn't exist.
	err = copyDir("/tmp/td/alderaan", "/tmp/td/lala")
	assert.NotNil(t, err)

	// One does not simply walk into Mordor...
	err = copyDir("/tmp/td/good", "/tmp/td/mordor/good")
	msg := "open /tmp/td/mordor/good/test.txt: permission denied"
	assert.Equal(t, err.Error(), msg)

	// ... but you can replace Mordor :D
	err = copyDir("/tmp/td/good", "/tmp/td/mordor")
	assert.Nil(t, err)
	body, err = ioutil.ReadFile("/tmp/td/mordor/test.txt")
	assert.Nil(t, err)
	assert.Equal(t, string(body), "hello world lala")
}

func TestRequestUrl(t *testing.T) {
	config = &configuration{
		Server: "http://localhost:9999/",
		Token:  "1234",
	}

	path1 := requestUrl("lala", true)
	path2 := requestUrl("/lala", true)
	assert.Equal(t, path1, path2)
	assert.Equal(t, path1, "http://localhost:9999/lala?token=1234")

	path3 := requestUrl("lala", false)
	assert.Equal(t, path3, "http://localhost:9999/lala")
}

func TestGetResponse(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	res, err := getResponse("GET", "/lala", nil)
	assert.Nil(t, err)
	assert.Equal(t, res.StatusCode, 200)

	req := res.Request
	assert.Equal(t, req.Method, "GET")
	assert.Equal(t, req.URL.Path, "/lala")
	assert.Equal(t, req.URL.RawQuery, "token=1234")
	assert.Equal(t, req.Header.Get("Content-Type"), "application/json")
}
