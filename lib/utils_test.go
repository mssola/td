// Copyright (C) 2014-2017 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mssola/capture"
)

func TestHome(t *testing.T) {
	casa := os.Getenv("HOME")
	defer func() { _ = os.Setenv("HOME", casa) }()

	if h := home(); casa != h {
		t.Fatalf("Expected '%v', Got '%v'", casa, h)
	}

	startTestEnv(t)
	defer stopTestEnv(t)

	errCheck(t, os.Setenv("TD", "/home/lala"))
	if h := home(); h != "/home/lala" {
		t.Fatalf("Expected '/home/lala', Got '%v'", h)
	}

	defer func() {
		if err := recover(); err != nil {
			msg := "You don't have the $HOME environment variable set"
			if err.(string) != msg {
				t.Fatalf("Expected '%v'; got: %v", msg, err)
			}
		}
	}()
	errCheck(t, os.Setenv("HOME", ""))
	errCheck(t, os.Setenv("TD", ""))
	home()
}

func TestEditor(t *testing.T) {
	errCheck(t, os.Setenv("EDITOR", "emacs"))
	if e := editor(); e != "emacs" {
		t.Fatalf("Expected 'emacs', Got '%v'", e)
	}

	errCheck(t, os.Setenv("EDITOR", ""))
	if e := editor(); e != defaultEditor {
		t.Fatalf("Expected '%v', Got '%v'", defaultEditor, e)
	}
}

func TestEditorArguments(t *testing.T) {
	defer func() { File = "" }()

	File = ""
	args := editorArguments()
	if len(args) != 0 {
		t.Fatalf("Expecting just 1 argument, got %v", len(args))
	}

	errCheck(t, os.Setenv("EDITOR", "emacs"))
	args = editorArguments()
	if len(args) != 0 {
		t.Fatalf("Expecting just 1 argument, got %v", len(args))
	}

	File = "something"
	res := capture.All(func() { args = editorArguments() })
	if len(args) != 0 {
		t.Fatalf("Expecting just 1 argument, got %v", len(args))
	}
	str := "the -f/--file flag does not work if it's not Vim"
	if !strings.Contains(string(res.Stdout), str) {
		t.Fatalf("Expecting '%s'; got: %s", str, res.Stdout)
	}

	errCheck(t, os.Setenv("EDITOR", "vim"))
	res = capture.All(func() { args = editorArguments() })
	if len(args) != 0 {
		t.Fatalf("Expecting just 1 argument, got %v", len(args))
	}
	str = "given file 'something' does not exist"
	if !strings.Contains(string(res.Stdout), str) {
		t.Fatalf("Expecting '%s'; got: %s", str, res.Stdout)
	}

	p, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not get working directory")
	}
	if strings.HasSuffix(p, "lib") {
		p = filepath.Dir(p)
	}
	File = filepath.Join(p, "README.md")

	args = editorArguments()
	if len(args) != 2 {
		t.Fatalf("Expecting 2 argument, got %v", len(args))
	}
	if args[0] != "-s" {
		t.Fatalf("Expecting '-s' argument, got %v", args[0])
	}
	if args[1] != File {
		t.Fatalf("Expecting '%v', argument, got %v", File, args[0])
	}
}

func TestCopyFile(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	// Prepare the filesystem.
	errCheck(t, os.RemoveAll("/tmp/td"))
	errCheck(t, os.MkdirAll("/tmp/td/good", 0755))
	errCheck(t, os.MkdirAll("/tmp/td/mordor", 0000))
	f, err := os.Create("/tmp/td/good/test.txt")
	defer func() { _ = f.Close() }()
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	_, err = f.WriteString("hello world")
	errCheck(t, err)

	// Good copy of the "good" dir.
	err = copyDir("/tmp/td/good", "/tmp/td/lala")
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	body, err := ioutil.ReadFile("/tmp/td/lala/test.txt")
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	if string(body) != "hello world" {
		t.Fatalf("Expected %v; got %v", string(body), "hello world")
	}

	// Copying the same directory replaces the old one.
	_, err = f.WriteString(" lala")
	errCheck(t, err)
	err = copyDir("/tmp/td/good", "/tmp/td/lala")
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	body, err = ioutil.ReadFile("/tmp/td/lala/test.txt")
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	if string(body) != "hello world lala" {
		t.Fatalf("Expected %v; got %v", string(body), "hello world lala")
	}

	// Trying from a directory that doesn't exist.
	err = copyDir("/tmp/td/alderaan", "/tmp/td/lala")
	if err == nil {
		t.Fatalf("We actually expected an error to happen here")
	}

	// One does not simply walk into Mordor...
	err = copyDir("/tmp/td/good", "/tmp/td/mordor/good")
	msg := "open /tmp/td/mordor/good/test.txt: permission denied"
	if err.Error() != msg {
		t.Fatalf("Expected %v; got %v", err.Error(), msg)
	}

	// ... but you can replace Mordor :D
	err = copyDir("/tmp/td/good", "/tmp/td/mordor")
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	body, err = ioutil.ReadFile("/tmp/td/mordor/test.txt")
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	if string(body) != "hello world lala" {
		t.Fatalf("Expected %v; got %v", string(body), "hello world lala")
	}
}

func TestRequestURL(t *testing.T) {
	Insecure = false
	config = &configuration{
		Server: "http://localhost:9999/",
		Token:  "1234",
	}
	_, err := requestURL("lala", true)
	msg := "attempted to reach a server that is not using HTTPS"
	if err.Error() != msg {
		t.Fatalf("Got: %v; Expected: %v", err.Error(), msg)
	}

	Insecure = true
	config = &configuration{
		Server: "http://localhost:9999/",
		Token:  "1234",
	}

	path1, err := requestURL("lala", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	path2, err := requestURL("/lala", true)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if path1 != path2 {
		t.Fatalf("'%v' should be the same as '%v'", path1, path2)
	}
	expected := "http://localhost:9999/lala?token=1234"
	if path1 != expected {
		t.Fatalf("'%v' should be the same as '%v'", path1, expected)
	}

	path3, err := requestURL("lala", false)
	expected = "http://localhost:9999/lala"
	if path3 != expected {
		t.Fatalf("'%v' should be the same as '%v'", path3, expected)
	}
}

func TestGetResponse(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	res, err := getResponse("GET", "/lala", nil)
	if err != nil {
		t.Fatalf("Did not expect error: %v", err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected %v; got %v", res.StatusCode, 200)
	}

	req := res.Request
	if req.Method != "GET" {
		t.Fatalf("Expected %v; got %v", req.Method, "GET")
	}
	if req.URL.Path != "/lala" {
		t.Fatalf("Expected %v; got %v", req.URL.Path, "/lala")
	}
	if req.URL.RawQuery != "token=1234" {
		t.Fatalf("Expected %v; got %v", req.URL.RawQuery, "token=1234")
	}
	if req.Header.Get("Content-Type") != "application/json" {
		t.Fatalf("Expected %v; got %v", req.Header.Get("Content-Type"), "application/json")
	}
}

func TestTimedOutGetResponse(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	// The trick to make it time out
	oldTimeout := requestTimeout
	requestTimeout = 500 * time.Millisecond
	defer func() { requestTimeout = oldTimeout }()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1 * time.Second)
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	_, err := getResponse("GET", "/lala", nil)
	if err == nil {
		t.Fatalf("Should've failed!")
	}
	if !strings.Contains(err.Error(), "timed out!") {
		t.Fatalf("Should've been a time out error, instead: %v", err)
	}
}

func TestErrorOnGetResponse(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrong 301 format, it should fail.
		http.Error(w, "some error", 301)
	}))
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	_, err := getResponse("GET", "/lala", nil)
	if err == nil {
		t.Fatalf("Should've failed!")
	}
	msg := "missing Location"
	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("Expected '%v'; got: %v", msg, err)
	}
}
