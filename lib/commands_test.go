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
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/mssola/capture"
)

var testTopics []Topic

// The parameters that be given through a request body.
type params struct {
	Name     string
	Contents string
}

type testOptions struct {
	Timeout     bool
	BadResponse bool
}

// Get the possible parameters from the given request. Note that it will only
// check for the "name" and "contents" parameters.
func getFromBody(req *http.Request) *params {
	var p params

	if req.Body == nil {
		return nil
	}

	decoder := json.NewDecoder(req.Body)
	if err := decoder.Decode(&p); err != nil {
		return nil
	}
	return &p
}

func compareSlices(t *testing.T, given, expected []string) {
	if len(given) != len(expected) {
		t.Fatalf("Wrong length between %#v -- %#v", given, expected)
	}
	for k, v := range given {
		if v != expected[k] {
			t.Fatalf("Given %v, Expected: %v", v, expected[k])
		}
	}
}

func urlIs(r *http.Request, url string) bool {
	return strings.HasPrefix(r.URL.String(), url)
}

func topicServer(opts *testOptions) *httptest.Server {
	testTopics = []Topic{
		{ID: "1", Name: "topic1", Contents: "1111"},
		{ID: "2", Name: "topic2", Contents: "2222"},
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if opts != nil && opts.Timeout {
			time.Sleep(1 * time.Second)
		}
		if !urlIs(r, "/topics") {
			return
		}

		switch r.Method {
		case "GET":
			b, _ := json.Marshal(testTopics)
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprint(w, string(b))
		case "POST":
			p := getFromBody(r)
			t := Topic{ID: p.Name, Name: p.Name}
			testTopics = append(testTopics, t)

			if opts == nil || !opts.BadResponse {
				b, _ := json.Marshal(&t)
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, string(b))
			}
		case "PUT":
			p := getFromBody(r)
			id := strings.Split(r.URL.Path, "/")[2]
			idx := 0

			for k, v := range testTopics {
				if v.ID == id {
					if p.Name != "" {
						testTopics[k].Name = p.Name
					}
					testTopics[k].Contents = p.Contents
					idx = k
				}
			}

			if opts == nil || !opts.BadResponse {
				b, _ := json.Marshal(&testTopics[idx])
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprint(w, string(b))
			}
		case "DELETE":
			parts := strings.Split(r.URL.Path, "/")

			b := testTopics[:0]
			length := 0
			for _, v := range testTopics {
				if v.ID != parts[2] {
					length++
					b = append(b, v)
				}
			}
			testTopics = testTopics[:length]

			if opts == nil || !opts.BadResponse {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(200)
			}
		}
	}))
}

func testList(t *testing.T, expected []string) {
	var err error
	res := capture.All(func() { err = List() })
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}

	output := strings.Split(strings.TrimSpace(string(res.Stdout)), "\n")
	compareSlices(t, output, expected)
}

func TestList(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(nil)
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	testList(t, []string{
		"Fetching the topics from the server.",
		"topic1",
		"topic2",
	})
}

func TestCreate(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(nil)
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	err := Create("topic3")
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	testList(t, []string{
		"Fetching the topics from the server.",
		"topic1",
		"topic2",
		"topic3",
	})
}

func TestBadServerCreate(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(&testOptions{BadResponse: true})
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	err := Create("topic3")
	if err == nil {
		t.Fatalf("We were expecting an error!")
	}
	msg := "could not create this topic"
	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("Expecting %v; Got: %v", msg, err.Error())
	}
}

func TestTimeoutCreate(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	oldTimeout := requestTimeout
	defer func() { requestTimeout = oldTimeout }()
	requestTimeout = 200 * time.Millisecond

	ts := topicServer(&testOptions{Timeout: true})
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	err := Create("topic3")
	if err == nil {
		t.Fatalf("We were expecting an error!")
	}
	msg := "timed out!"
	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("Expecting %v; Got: %v", msg, err.Error())
	}
}

func TestDelete(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(nil)
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	var err error
	capture.All(func() { err = fetch() })
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	if err = Delete("topic1"); err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	testList(t, []string{
		"Fetching the topics from the server.",
		"topic2",
	})
}

func TestUnknownTopicDelete(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(nil)
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	var err error
	capture.All(func() { err = fetch() })
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	if err = Delete("topics1"); err == nil {
		t.Fatal("We were expecting an error")
	}

	parts := strings.Split(err.Error(), "\n")
	if !strings.Contains(parts[0], "the topic 'topics1' does not exist") {
		t.Fatalf("Not expecting: %v", parts[0])
	}
	if !strings.Contains(parts[2], "Did you mean one of these?") {
		t.Fatalf("Not expecting: %v", parts[2])
	}
	if !strings.Contains(parts[3], "topic1") {
		t.Fatalf("Not expecting: %v", parts[3])
	}
	if !strings.Contains(parts[4], "topic2") {
		t.Fatalf("Not expecting: %v", parts[4])
	}
}

func TestTimeoutDelete(t *testing.T) {
	// This does not work in go 1.5 for some reason :/
	if strings.Contains(runtime.Version(), "1.5") {
		return
	}

	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(&testOptions{Timeout: true})
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	var err error
	capture.All(func() { err = fetch() })
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}

	oldTimeout := requestTimeout
	defer func() { requestTimeout = oldTimeout }()
	requestTimeout = 10 * time.Millisecond

	// TODO: this doesn't work for some reason :/
	/*
		if err = Delete("topic1"); err == nil {
			t.Fatal("We were expecting an error")
		}
		if !strings.Contains(err.Error(), "timed out") {
			t.Fatal("We were expecting a time out error")
		}
	*/
}

func TestRename(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(nil)
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	var err error
	capture.All(func() { err = fetch() })
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	if err = Rename("topic1", "newtopic"); err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	testList(t, []string{
		"Fetching the topics from the server.",
		"newtopic",
		"topic2",
	})
}

func TestUnknownTopicRename(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(nil)
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	var err error
	capture.All(func() { err = fetch() })
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	if err = Rename("1", "topic"); err == nil {
		t.Fatal("We were expecting an error")
	}

	parts := strings.Split(err.Error(), "\n")
	if !strings.Contains(parts[0], "the topic '1' does not exist") {
		t.Fatalf("Not expecting: %v", parts[0])
	}
	if len(parts) != 2 {
		t.Fatalf("Expected length: 2; got: %v", len(parts))
	}
}

func TestTimeoutRename(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(&testOptions{Timeout: true})
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	var err error
	capture.All(func() { err = fetch() })
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}

	oldTimeout := requestTimeout
	defer func() { requestTimeout = oldTimeout }()
	requestTimeout = 200 * time.Millisecond

	if err = Rename("topic1", "topic2"); err == nil {
		t.Fatal("We were expecting an error")
	}
	if !strings.Contains(err.Error(), "timed out") {
		t.Fatal("We were expecting an error")
	}
}

func TestBadResponseRename(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(&testOptions{BadResponse: true})
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	var err error
	capture.All(func() { err = fetch() })
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}

	err = Rename("topic1", "topic")
	if err == nil {
		t.Fatalf("We were expecting an error!")
	}
	msg := "could not rename this topic"
	if !strings.Contains(err.Error(), msg) {
		t.Fatalf("Expecting %v; Got: %v", msg, err.Error())
	}
}

func TestEdit(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(nil)
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	oldCommand := editCommand
	defer func() { editCommand = oldCommand }()
	editCommand = func() error {
		path := filepath.Join(home(), dirName, newDir, "topic1.md")
		return ioutil.WriteFile(path, []byte("contents"), 0755)
	}

	if testTopics[0].Contents != "1111" {
		t.Fatalf("Expecting \"1111\"; got: %v", testTopics[0].Contents)
	}

	var err error
	capture.All(func() { err = Edit() })

	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}
	if testTopics[0].ID != "1" {
		t.Fatalf("Expecting \"1\"; got: %v", testTopics[0].ID)
	}
	if testTopics[0].Name != "topic1" {
		t.Fatalf("Expecting \"topic1\"; got: %v", testTopics[0].Name)
	}
	if testTopics[0].Contents != "contents" {
		t.Fatalf("Expecting \"contents\"; got: %v", testTopics[0].Contents)
	}
}

func TestEditCommand(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := topicServer(nil)
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	_ = os.Setenv("EDITOR", "echo")
	var err error
	res := capture.All(func() { err = Edit() })
	if err != nil {
		t.Fatalf("Not expecting error: %v", err)
	}
	parts := strings.Split(string(res.Stdout), "\n")
	compareSlices(t, parts, []string{
		"Fetching the topics from the server.",
		"", // echo
		"",
	})
}
