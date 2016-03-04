// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTopicsConfig(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)

	Initialize()

	// Reading an empty config file.
	var topics []Topic
	readTopics(&topics)
	if len(topics) != 0 {
		t.Fatalf("Expected length %v; got %v", 0, len(topics))
	}

	// Let's add a topic.
	topic1 := &Topic{ID: "1", Name: "topic1"}
	addTopic(topic1)
	readTopics(&topics)
	if len(topics) != 1 {
		t.Fatalf("Expected length %v; got %v", 1, len(topics))
	}
	if topic1.ID != topics[0].ID {
		t.Fatalf("Expected %v; got %v", topic1.ID, topics[0].ID)
	}

	// Writing two topics all at once (while replacing "topic1").
	var addedTopics []Topic
	addedTopics = append(addedTopics, Topic{ID: "2", Name: "topic2"})
	addedTopics = append(addedTopics, Topic{ID: "3", Name: "topic3"})
	writeTopics(addedTopics)
	readTopics(&topics)
	if len(topics) != 2 {
		t.Fatalf("Expected length %v; got %v", 2, len(topics))
	}
	if "2" != topics[0].ID {
		t.Fatalf("Expected 2; got %v", topics[0].ID)
	}
	if "3" != topics[1].ID {
		t.Fatalf("Expected 3; got %v", topics[1].ID)
	}
}

func pathFor(prefix, file string) string {
	return filepath.Join(os.Getenv("TD"), prefix, file)
}

func TestSave(t *testing.T) {
	startTestEnv(t)
	defer stopTestEnv(t)
	Initialize()

	// Save some topics.
	var addedTopics []Topic
	addedTopics = append(addedTopics, Topic{ID: "1", Name: "topic1"})
	addedTopics = append(addedTopics, Topic{ID: "2", Name: "topic2"})
	save(addedTopics)

	// Now let's see if the files are in order.
	tmpContents1, _ := ioutil.ReadFile(pathFor(tmpDir, "topic1.md"))
	tmpContents2, _ := ioutil.ReadFile(pathFor(tmpDir, "topic2.md"))
	newContents1, _ := ioutil.ReadFile(pathFor(newDir, "topic1.md"))
	newContents2, _ := ioutil.ReadFile(pathFor(newDir, "topic2.md"))
	oldContents1, _ := ioutil.ReadFile(pathFor(oldDir, "topic1.md"))
	oldContents2, _ := ioutil.ReadFile(pathFor(oldDir, "topic2.md"))

	// Comparing contents!
	if string(tmpContents1) != string(newContents1) {
		t.Fatalf("Expected %v; got %v", tmpContents1, newContents1)
	}
	if string(tmpContents1) != string(oldContents1) {
		t.Fatalf("Expected %v; got %v", tmpContents1, oldContents1)
	}
	if string(tmpContents2) != string(newContents2) {
		t.Fatalf("Expected %v; got %v", tmpContents2, newContents2)
	}
	if string(tmpContents2) != string(oldContents2) {
		t.Fatalf("Expected %v; got %v", tmpContents2, oldContents2)
	}
}
