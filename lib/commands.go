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
	"strings"
	"time"
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
	adir := filepath.Join(home(), dirName, newDir)
	entries, _ := ioutil.ReadDir(adir)
	for _, entry := range entries {
		name := strings.SplitN(entry.Name(), ".", 2)
		fmt.Printf("%v\n", name[0])
	}
	return nil
}

func Push() error {
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
	return nil
}

func Rename(oldName, newName string) error {
	return nil
}
