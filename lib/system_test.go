// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTopicsConfig(t *testing.T) {
	// Give a warm place to this test.
	os.RemoveAll("/tmp/td")
	os.MkdirAll("/tmp/td", 0755)
	os.Setenv("TD", "/tmp")
	dirName = "td"
	Initialize()
	_, err := os.Stat("/tmp/td/config.json")
	assert.Nil(t, err)

	// Reading an empty config file.
	var topics []Topic
	readTopics(&topics)
	assert.Equal(t, len(topics), 0)

	// Let's add a topic.
	topic1 := &Topic{Id: "1", Name: "topic1"}
	addTopic(topic1)
	readTopics(&topics)
	assert.Equal(t, len(topics), 1)
	assert.Equal(t, topic1.Id, topics[0].Id)

	// Writing two topics all at once (while replacing "topic1").
	var addedTopics []Topic
	addedTopics = append(addedTopics, Topic{Id: "2", Name: "topic2"})
	addedTopics = append(addedTopics, Topic{Id: "3", Name: "topic3"})
	writeTopics(addedTopics)
	readTopics(&topics)
	assert.Equal(t, len(topics), 2)
	assert.Equal(t, "2", topics[0].Id)
	assert.Equal(t, "3", topics[1].Id)

	// Tearing down.
	os.Setenv("TD", "")
	dirName = ".td"
}

func TestSave(t *testing.T) {
	// Give a warm place to this test.
	os.RemoveAll("/tmp/td")
	os.MkdirAll("/tmp/td", 0755)
	os.Setenv("TD", "/tmp")
	dirName = "td"
	Initialize()
	_, err := os.Stat("/tmp/td/config.json")
	assert.Nil(t, err)

	// Save some topics.
	var addedTopics []Topic
	addedTopics = append(addedTopics, Topic{Id: "1", Name: "topic1"})
	addedTopics = append(addedTopics, Topic{Id: "2", Name: "topic2"})
	save(addedTopics)

	// Now let's see if the files are in order.
	tmpContents1, _ := ioutil.ReadFile("/tmp/td/tmp/topic1.md")
	tmpContents2, _ := ioutil.ReadFile("/tmp/td/tmp/topic2.md")
	newContents1, _ := ioutil.ReadFile("/tmp/td/new/topic1.md")
	newContents2, _ := ioutil.ReadFile("/tmp/td/new/topic2.md")
	oldContents1, _ := ioutil.ReadFile("/tmp/td/old/topic1.md")
	oldContents2, _ := ioutil.ReadFile("/tmp/td/old/topic2.md")

	// Comparing contents!
	assert.Equal(t, tmpContents1, newContents1)
	assert.Equal(t, tmpContents1, oldContents1)
	assert.Equal(t, tmpContents2, newContents2)
	assert.Equal(t, tmpContents2, oldContents2)

	// Tearing down.
	os.Setenv("TD", "")
	dirName = ".td"
}

func TestUpdate(t *testing.T) {
	h := testHelper()
	defer h.teardown()

	// Give a warm place to this test.
	os.RemoveAll("/tmp/td")
	os.MkdirAll("/tmp/td", 0755)
	os.Setenv("TD", "/tmp")
	dirName = "td"
	Initialize()
	_, err := os.Stat("/tmp/td/config.json")
	assert.Nil(t, err)

	// Let's write something new to the "new" directory.
	f, _ := os.Create("/tmp/td/new/lala.md")
	f.WriteString("lala")
	f.Close()

	// This is successful, but does nothing at all.
	update([]string{}, []string{})
	oldContents, err := ioutil.ReadFile("/tmp/td/old/lala.md")
	assert.NotNil(t, err)

	// Fails do nothing to the FS.
	update([]string{}, []string{"lala"})
	oldContents, err = ioutil.ReadFile("/tmp/td/old/lala.md")
	assert.NotNil(t, err)

	// This is successful and copies the recently created "lala.md" file into
	// the "old" directory.
	update([]string{"lala"}, []string{})
	oldContents, err = ioutil.ReadFile("/tmp/td/old/lala.md")
	assert.Nil(t, err)
	newContents, _ := ioutil.ReadFile("/tmp/td/new/lala.md")
	assert.Equal(t, oldContents, newContents)

	// Tearing down.
	os.Setenv("TD", "")
	dirName = ".td"
}
