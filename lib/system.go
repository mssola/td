// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

var (
	dirName = ".td"
)

const (
	topicsName = "topics.json"
	tmpDir     = "tmp"
	oldDir     = "old"
	newDir     = "new"
)

func readTopics(topics *[]Topic) {
	file := filepath.Join(home(), dirName, topicsName)
	body, _ := ioutil.ReadFile(file)
	json.Unmarshal(body, topics)
}

func writeTopics(topics []Topic) {
	// Clean it up, we don't want to store the contents.
	for k, _ := range topics {
		topics[k].Contents = ""
		topics[k].Markdown = ""
	}
	body, _ := json.Marshal(topics)

	// Write the JSON.
	file := filepath.Join(home(), dirName, topicsName)
	f, _ := os.Create(file)
	f.Write(body)
	f.Close()
}

func addTopic(topic *Topic) {
	// Add the topic to the JSON file.
	topics := []Topic{*topic}
	readTopics(&topics)
	writeTopics(topics)

	// And create the files for this new topic.
	odir := filepath.Join(home(), dirName, oldDir)
	write(topic, odir)
	odir = filepath.Join(home(), dirName, newDir)
	write(topic, odir)
}

func save(topics []Topic) {
	// First of all, reset the temporary directory.
	dir := filepath.Join(home(), dirName, tmpDir)
	os.RemoveAll(dir)
	os.MkdirAll(dir, 0755)

	// Save all the topics to this temporary directory.
	for _, t := range topics {
		write(&t, dir)
	}

	// Update the old and new directories
	adir := filepath.Join(home(), dirName, oldDir)
	copyDir(dir, adir)
	adir = filepath.Join(home(), dirName, newDir)
	copyDir(dir, adir)

	// And finally, write the JSON file.
	writeTopics(topics)
}

func write(topic *Topic, path string) {
	path = filepath.Join(path, topic.Name+".md")
	f, _ := os.Create(path)
	f.Close()
	f.WriteString(topic.Contents)
}

func update(success, fails []string) {
	srcDir := filepath.Join(home(), dirName, newDir)
	dstDir := filepath.Join(home(), dirName, oldDir)

	// Copy successes.
	for _, v := range success {
		src := filepath.Join(srcDir, v+".md")
		dst := filepath.Join(dstDir, v+".md")
		copyFile(src, dst)
	}

	// List failures.
	if len(fails) == 0 {
		fmt.Printf("Success!\n")
	} else {
		fmt.Printf("The following topics could not be pushed:\n")
		for _, v := range fails {
			fmt.Printf("\t" + v + "\n")
		}
	}
}
