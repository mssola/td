// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"time"

	"github.com/mssola/dym"
)

// Topic contains all the information related to a topic.
type Topic struct {
	ID        string    `json:"id,omitempty"`
	Name      string    `json:"name,omitempty"`
	Contents  string    `json:"contents,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	Markdown  string    `json:"markdown,omitempty"`
	Error     string    `json:"error,omitempty"`
}

// unknownTopic returns the proper error in the case that the given topic
// does not exist.
func unknownTopic(name string) error {
	var topics []Topic
	var names []string

	readTopics(&topics)
	for _, v := range topics {
		names = append(names, v.Name)
	}

	msg := fmt.Sprintf("the topic '%v' does not exist", name)
	similars := dym.Similar(names, name)
	if len(similars) == 0 {
		return NewError(msg)
	}

	msg += "\n\nDid you mean one of these?\n"
	for _, v := range similars {
		msg += "\t" + v + "\n"
	}
	return NewError(msg)
}

// topicResponse parses the given response and fill the given topic with the
// abstracted information.
func topicResponse(t *Topic, res *http.Response) error {
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, t); err != nil {
		return errors.New("unknown topic format")
	}
	if t.Error != "" {
		return errors.New(t.Error)
	}
	return nil
}

// safeFetch returns whether it's safe to fetch topics from the server or not.
// This depends on whether there are changes that have not been pushed or not.
func safeFetch() bool {
	topics := changedTopics()
	return len(topics) == 0
}

// fetch saves all the topics from the server locally.
func fetch() error {
	// Ask before wiping out the currently changed topics.
	if !safeFetch() {
		return errors.New("you have changes on the currently cached topics")
	}

	// Perform the HTTP request.
	fmt.Printf("Fetching the topics from the server.\n")
	res, err := getResponse("GET", "/topics", nil)
	if err != nil {
		return err
	}

	// Parse the given topics.
	var topics []Topic
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &topics); err != nil {
		return fromError(err)
	}

	// And save the results.
	save(topics)
	return nil
}

// pushTopics pushes all the given topics to the server. Only successful pushes
// will be updated locally.
func pushTopics(topics []Topic) {
	var success, fails []string

	total := len(topics)
	for k, v := range topics {
		// Print the status.
		fmt.Printf("\rPushing... %v/%v\r", k+1, total)

		// Get the contents.
		file := filepath.Join(home(), dirName, newDir, v.Name+".md")
		body, _ := ioutil.ReadFile(file)
		t := &Topic{Contents: string(body)}
		if t.Contents == "" {
			success = append(success, v.Name)
			continue
		}

		// Perform the request.
		body, _ = json.Marshal(t)
		path := "/topics/" + v.ID
		_, err := getResponse("PUT", path, bytes.NewReader(body))
		if err == nil {
			success = append(success, v.Name)
		} else {
			fails = append(fails, v.Name)
		}
	}

	// And finally update the file system.
	update(success, fails)
}
