// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// This is the subdirectory inside "home" where all the application data is
	// contained in.
	dirName = ".td"
)

const (
	// The name of the list of topics.
	topicsName = "topics.json"

	// The name for the directory where temporary data gets stored.
	tmpDir = "tmp"

	// The name for the directory with contents from the latest version on the
	// server.
	oldDir = "old"

	// The name for the directory where the user edits topics (a.k.a. the "To
	// do" list).
	newDir = "new"
)

// Note that all these functions never return an error. This is because all the
// errors that could be returned are I/O related, and any error of this kind
// has already been checked because of the initial call to the "Initialize"
// function in the "main" function. Therefore, it's not worth to check errors,
// since they will, in fact, never occur.

// Read all the topics that we have localy and put them in the given topics
// array.
func readTopics(topics *[]Topic) {
	file := filepath.Join(home(), dirName, topicsName)
	body, _ := ioutil.ReadFile(file)
	json.Unmarshal(body, topics)
}

// Save the given topics into the list of local topics. Note that this function
// effectively replaces the previous list.
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

// Add the given topic to the list of local topics.
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

// Returns a list of all the topics that have changed since the last version.
// TODO: test me!
func changedTopics() []Topic {
	var topics, changed []Topic
	readTopics(&topics)

	re, _ := regexp.Compile(".*/(.+)\\.md")
	sDir := filepath.Join(home(), dirName, oldDir)
	dDir := filepath.Join(home(), dirName, newDir)
	out, _ := exec.Command("diff", "-qr", sDir, dDir).Output()

	for _, l := range strings.Split(string(out), "\n") {
		// Filter out blank lines.
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		// And now append the topic.
		match := re.FindSubmatch([]byte(l))
		if match != nil && len(match) == 2 {
			name := string(match[1])
			for _, v := range topics {
				if v.Name == name {
					changed = append(changed, v)
				}
			}
		}
	}
	return changed
}

// Save all the data from the given topics. This means that all the directories
// will be updates accordingly with the new contents for each file. Plus, this
// function will also call the "writeTopics" function in order to store the
// given list of topics into our local list of topics.
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

// Save the contents of the given topic. The file getting created will be the
// name of the topic with the ".md" extension. The directory where this file
// will be contained is the given "path" parameter.
func write(topic *Topic, path string) {
	path = filepath.Join(path, topic.Name+".md")
	f, _ := os.Create(path)
	f.WriteString(topic.Contents)
	f.Close()
}

// Copy all the files from the "new" directory to the "old" directory. This is
// done when performing the "push" command. Also note that this function will
// print the list of topics that could not be pushed if any. The two given
// parameters are a list of strings containing names of topics. The "success"
// slice contains the topics that have already been pushed to the server. The
// "fails" slice contains the topics that have failed on the push action.
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
