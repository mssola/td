// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/mssola/dym"
)

type Topic struct {
	Id         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Contents   string    `json:"contents,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
	Markdown   string    `json:"markdown,omitempty"`
	Error      string    `json:"error,omitempty"`
}

func unknownTopic(name string) {
	var topics []Topic
	var names []string

	readTopics(&topics)
	for _, v := range topics {
		names = append(names, v.Name)
	}

	msg := fmt.Sprintf("td: the topic '%v' does not exist.", name)
	similars := dym.Similar(names, name)
	if len(similars) == 0 {
		fmt.Printf(msg)
	} else {
		msg += "\n\nDid you mean one of these?\n"
		for _, v := range similars {
			msg += "\t" + v + "\n"
		}
		fmt.Printf(msg)
	}
}

func Edit() error {
	cmd := exec.Command(editor())
	cmd.Dir = filepath.Join(home(), dirName, newDir)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func Fetch() error {
	// Perform the HTTP request.
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
	fmt.Printf("Topics updated.\n")
	return nil
}

func List() error {
	var topics []Topic
	readTopics(&topics)
	for _, v := range topics {
		fmt.Printf("%v\n", v.Name)
	}
	return nil
}

func Push() error {
	var success, fails []string
	var topics []Topic
	readTopics(&topics)

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
		path := "/topics/" + v.Id
		_, err := getResponse("PUT", path, bytes.NewReader(body))
		if err == nil {
			success = append(success, v.Name)
		} else {
			fails = append(fails, v.Name)
		}
	}

	// And finally update the file system.
	update(success, fails)
	return nil
}

func Status() error {
	files := []string{}
	re, _ := regexp.Compile("(\\w+)\\.md")

	sDir := filepath.Join(home(), dirName, oldDir)
	dDir := filepath.Join(home(), dirName, newDir)
	out, _ := exec.Command("diff", "-qr", sDir, dDir).Output()

	for _, l := range strings.Split(string(out), "\n") {
		// Filter out blank lines.
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}

		match := re.FindSubmatch([]byte(l))
		if match != nil && len(match) == 2 {
			files = append(files, string(match[1]))
		}
	}

	if len(files) > 0 {
		fmt.Printf("The following topics have changed since the last " +
			"version:\n\n")
		for _, f := range files {
			fmt.Printf("\t%v\n", f)
		}
	} else {
		fmt.Printf("There's nothing to be pushed.\n")
	}
	return nil
}

func topicResponse(t *Topic, res *http.Response) bool {
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, t); err != nil {
		return false
	}
	if t.Error != "" {
		return false
	}
	return true
}

func Create(name string) error {
	// Perform the HTTP request.
	t := &Topic{Name: name}
	body, _ := json.Marshal(t)
	res, err := getResponse("POST", "/topics", bytes.NewReader(body))
	if err != nil {
		return err
	}

	// Parse the newly created topic and add it to the list.
	if !topicResponse(t, res) {
		return newError("could not create this topic")
	}
	addTopic(t)
	return nil
}

func Delete(name string) error {
	var topics, actual []Topic
	var id string

	// Get the list of topics straight.
	readTopics(&topics)
	for _, v := range topics {
		if v.Name == name {
			id = v.Id
		} else {
			actual = append(actual, v)
		}
	}
	if id == "" {
		unknownTopic(name)
		os.Exit(1)
	}

	// Perform the HTTP request.
	if _, err := getResponse("DELETE", "/topics/"+id, nil); err != nil {
		return err
	}

	// On the system.
	writeTopics(actual)
	file := filepath.Join(home(), dirName, oldDir, name+".md")
	os.RemoveAll(file)
	file = filepath.Join(home(), dirName, newDir, name+".md")
	os.RemoveAll(file)
	return nil
}

func Rename(oldName, newName string) error {
	var topics []Topic
	var id, name string

	readTopics(&topics)
	for k, v := range topics {
		if v.Name == oldName {
			id = v.Id
			topics[k].Name = newName
		}
	}
	if id == "" {
		unknownTopic(name)
		os.Exit(1)
	}

	// Perform the HTTP Request.
	t := &Topic{Name: newName}
	body, _ := json.Marshal(t)
	res, err := getResponse("PUT", "/topics/"+id, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if !topicResponse(t, res) {
		return newError("could not rename this topic")
	}

	// Update the system.
	writeTopics(topics)
	file := filepath.Join(home(), dirName, oldDir)
	os.Rename(filepath.Join(file, oldName), filepath.Join(file, newName))
	file = filepath.Join(home(), dirName, newDir)
	os.Rename(filepath.Join(file, oldName), filepath.Join(file, newName))
	return nil
}
