// Copyright (C) 2014-2015 Miquel Sabaté Solà <mikisabate@gmail.com>
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
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	// The fallback editor in case the $EDITOR environment variable is not set.
	defaultEditor = "vi"
)

var (
	// The timeout for any HTTP request.
	requestTimeout = 15 * time.Second
)

// Returns the value of the current home. This value is fetched from the $TD
// environment variable. If it's not set, then the $HOME environment variable
// will be picked. If the $HOME environment variable is not set either,
// then it panics.
func home() string {
	value := os.Getenv("TD")
	if value == "" {
		value = os.Getenv("HOME")
		if value == "" {
			panic("You don't have the $HOME environment variable set")
		}
	}
	return value
}

// Returns the value of the $EDITOR environment variable. If this variable is
// not set, then it will return the value of the "defaultEditor" constant.
func editor() string {
	value := os.Getenv("EDITOR")
	if value == "" {
		return defaultEditor
	}
	return value
}

// Copy a file from a source path to a destination path. This function assumes
// that the source path exists. The only error that can be tolerated is if the
// user is trying to cpy a file into a protected directory.
func copyFile(source string, dest string) error {
	sf, _ := os.Open(source)
	defer sf.Close()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer df.Close()
	_, err = io.Copy(df, sf)
	return err
}

// Copy a given source directory into a destination directory. If the
// destination directory already exists, it will be removed. This function will
// only tolerate the following errors:
//
//	1. The copying of the files inside a directory has failed.
//  2. The source directory cannot be read.
//
// Note that subdirectories will *not* be copied. This is because this is a
// simple function that adjust to our scheme: directories only have regular
// files inside.
func copyDir(source string, dest string) error {
	os.RemoveAll(dest)
	os.MkdirAll(dest, 0755)

	entries, err := ioutil.ReadDir(source)
	if err != nil {
		return fromError(err)
	}

	for _, entry := range entries {
		sfp := filepath.Join(source, entry.Name())
		dfp := filepath.Join(dest, entry.Name())
		err = copyFile(sfp, dfp)
		if err != nil {
			return err
		}
	}
	return nil
}

// Construct the URL for the given path. The second parameter "token" tells
// this function whether it should include the authorization token in the
// query.
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

// Perform an HTTP request and get back the response. The "method" parameter
// corresponds to an HTTP method (e.g. "GET") and the "url" parameter corresponds
// to just the path for the URL (e.g. "/topics"). Some HTTP requests might want
// to send data through the body of the request. In this case the "body"
// parameter should be used.
func getResponse(method, url string, body io.Reader) (*http.Response, error) {
	// Setup the client for the HTTP request.
	client := http.Client{
		Timeout: requestTimeout,
	}
	str := requestUrl(url, true)
	req, _ := http.NewRequest(method, str, body)
	req.Header.Set("Content-Type", "application/json")

	// Perform the HTTP request itself.
	res, err := client.Do(req)
	if err == nil {
		return res, nil
	}

	// Check specifically for a timeout.
	if strings.HasSuffix(err.Error(), "use of closed network connection") {
		return nil, newError("timed out! Try it again in another time")
	}

	// Beautify the given error: only return the actual message.
	re, _ := regexp.Compile(`:\s+(.+)$`)
	errString := []byte(err.Error())
	if e := re.FindSubmatch(errString); len(e) == 2 {
		return nil, newError(string(e[1]))
	}
	return nil, fromError(err)
}

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
}

// TODO
func safeFetch() bool {
	if topics := changedTopics(); len(topics) > 0 {
		fmt.Printf("You haven't pushed the changes for the following topics:\n")
		for _, v := range topics {
			fmt.Printf("\t%v\n", v)
		}
		fmt.Printf("Are you sure that you want to do this? (y/n): ")

		var str string
		fmt.Scanln(&str)
		str = strings.ToLower(str)
		return str == "y" || str == "yes"
	}
	return true
}
