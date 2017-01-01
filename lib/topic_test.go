// Copyright (C) 2014-2017 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBadTopicResponse(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "h")
	}))
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	res, err := getResponse("POST", "/topics", nil)
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	topic := &Topic{}
	err = topicResponse(topic, res)
	if err == nil || err.Error() != "unknown topic format" {
		t.Fatalf("Expected error")
	}
}

func TestErrorResponse(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(&Topic{Error: "error"})
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, string(b))
	}))
	defer ts.Close()

	config = &configuration{
		Server: ts.URL,
		Token:  "1234",
	}

	res, err := getResponse("POST", "/topics", nil)
	if err != nil {
		t.Fatalf("We were not expecting an error: %v", err)
	}
	topic := &Topic{}
	err = topicResponse(topic, res)
	if err == nil || err.Error() != "error" {
		t.Fatalf("Expected error")
	}
}
