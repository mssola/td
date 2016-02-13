// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Done this way to test it.
var editCommand = func() error {
	cmd := exec.Command(editor())
	cmd.Dir = filepath.Join(home(), dirName, newDir)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// Edit performs the default command. That is, it fetches all the topics, opens
// up the default editor and pushes the changes.
func Edit() error {
	// Fetch the topics from the server.
	if err := fetch(); err != nil {
		return err
	}

	// Open up the editor.
	if err := editCommand(); err != nil {
		return fromError(err)
	}

	// Push all the changed files.
	changed := changedTopics()
	if len(changed) > 0 {
		fmt.Printf("Pushing your changes to the server.\n")
		pushTopics(changed)
	}
	return nil
}

// List simply shows the currently available topics.
func List() error {
	// Try to fetch them if no one else has done it. We can safely ignore the
	// error since we can still cache it if it exists. Otherwise it's not such
	// a pain to get an empty list on weird scenarios.
	_ = fetch()

	var topics []Topic
	readTopics(&topics)
	for _, v := range topics {
		fmt.Printf("%v\n", v.Name)
	}
	return nil
}

// Create creates a new topic on the server.
func Create(name string) error {
	// Perform the HTTP request.
	t := &Topic{Name: name}
	body, _ := json.Marshal(t)
	res, err := getResponse("POST", "/topics", bytes.NewReader(body))
	if err != nil {
		return err
	}

	// Parse the newly created topic and add it to the list.
	if err = topicResponse(t, res); err != nil {
		return NewError("could not create this topic")
	}
	addTopic(t)
	return nil
}

// Delete deletes the specified topic from the server.
func Delete(name string) error {
	var topics, actual []Topic
	var id string

	// Get the list of topics straight.
	readTopics(&topics)
	for _, v := range topics {
		if v.Name == name {
			id = v.ID
		} else {
			actual = append(actual, v)
		}
	}
	if id == "" {
		return unknownTopic(name)
	}

	// Perform the HTTP request.
	if _, err := getResponse("DELETE", "/topics/"+id, nil); err != nil {
		return err
	}

	// On the system.
	writeTopics(actual)
	file := filepath.Join(home(), dirName, oldDir, name+".md")
	_ = os.RemoveAll(file)
	file = filepath.Join(home(), dirName, newDir, name+".md")
	_ = os.RemoveAll(file)
	return nil
}

// Rename changes the name of the given topic with the new one.
func Rename(oldName, newName string) error {
	var topics []Topic
	var id string

	readTopics(&topics)
	for k, v := range topics {
		if v.Name == oldName {
			id = v.ID
			topics[k].Name = newName
		}
	}
	if id == "" {
		return unknownTopic(oldName)
	}

	// Perform the HTTP Request.
	t := &Topic{Name: newName}
	body, _ := json.Marshal(t)
	res, err := getResponse("PUT", "/topics/"+id, bytes.NewReader(body))
	if err != nil {
		return err
	}
	if err = topicResponse(t, res); err != nil {
		return NewError("could not rename this topic")
	}

	// Update the system.
	writeTopics(topics)
	file := filepath.Join(home(), dirName, oldDir)
	_ = os.Rename(filepath.Join(file, oldName), filepath.Join(file, newName))
	file = filepath.Join(home(), dirName, newDir)
	_ = os.Rename(filepath.Join(file, oldName), filepath.Join(file, newName))
	return nil
}
