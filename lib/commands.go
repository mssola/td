// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/mssola/dym"
)

type Topic struct {
	Id         string    `json:"id,omitempty"`
	Name       string    `json:"name,omitempty"`
	Contents   string    `json:"contents,omitempty"`
	Created_at time.Time `json:"created_at,omitempty"`
	Markdown   string    `json:"markdown,omitempty"`
}

// TODO
func requestUrl(path string, token bool) string {
	u, _ := url.Parse(config.Server)
	u.Path = path
	if token {
		v := url.Values{}
		v.Set("token", config.Token)
		u.RawQuery = v.Encode()
	}
	return u.String()
}

func getResponse(method, url string, body io.Reader) (*http.Response, error) {
	client := &http.Client{}
	str := requestUrl(url, true)
	req, _ := http.NewRequest(method, str, body)
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func unknownTopic(name string) {
	var topics []Topic
	var names []string

	getTopics(&topics)
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
		return fromError(err)
	}

	// Parse the given topics.
	var topics []Topic
	body, _ := ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &topics); err != nil {
		return fromError(err)
	}

	// And save the results.
	if err := save(topics); err != nil {
		return err
	}
	fmt.Printf("Topics updated.\n")
	return nil
}

func List() error {
	var topics []Topic
	getTopics(&topics)
	for _, v := range topics {
		fmt.Printf("%v\n", v.Name)
	}
	return nil
}

func Push() error {
	var success, fails []string
	var topics []Topic
	getTopics(&topics)

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

	update(success, fails)
	fmt.Printf("Success!\n")
	return nil
}

func Create(name string) error {
	// Perform the HTTP request.
	t := &Topic{Name: name}
	body, _ := json.Marshal(t)
	res, err := getResponse("POST", "/topics", bytes.NewReader(body))
	if err != nil {
		return fromError(err)
	}

	// Parse the newly created topic and add it to the list.
	body, _ = ioutil.ReadAll(res.Body)
	if err := json.Unmarshal(body, &t); err != nil {
		return fromError(err)
	}
	return addTopic(t)
}

func Delete(name string) error {
	var topics, actual []Topic
	var tid, tname string

	// Get the list of topics straight.
	getTopics(&topics)
	for _, v := range topics {
		if v.Name == name {
			tid, tname = v.Id, v.Name
		} else {
			actual = append(actual, v)
		}
	}
	if tid == "" {
		unknownTopic(name)
		os.Exit(1)
	}

	// Perform the HTTP request.
	_, err := getResponse("DELETE", "/topics/"+tid, nil)
	if err != nil {
		return fromError(err)
	}

	// On the system.
	writeJson(actual)
	file := filepath.Join(home(), dirName, oldDir, tname+".md")
	os.RemoveAll(file)
	file = filepath.Join(home(), dirName, newDir, tname+".md")
	os.RemoveAll(file)
	return nil
}

func Rename(oldName, newName string) error {
	return nil
}
