// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"
)

type Topic struct {
	Id         string    `json:"id"`
	Name       string    `json:"name"`
	Contents   string    `json:"contents"`
	Created_at time.Time `json:"created_at"`
	Markdown   string    `json:"markdown"`
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
	return nil
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
		parts := strings.SplitN(entry.Name(), "_", 2)
		name := strings.SplitN(parts[1], ".", 2)
		fmt.Printf("%v\n", name[0])
	}
	return nil
}

func Push() error {
	return nil
}

func Create(name string) error {
	return nil
}

func Delete(name string) error {
	return nil
}

func Rename(oldName, newName string) error {
	return nil
}
