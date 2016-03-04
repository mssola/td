// Copyright (C) 2014-2016 Miquel Sabaté Solà <mikisabate@gmail.com>
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package lib

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net"
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

	// Insecure contains whether HTTP communications are allowed instead of
	// secure HTTPS ones. Defaults to false.
	Insecure = false

	// TLSVerify sets whether certificates have to be validated. Defaults to
	// true. Ignored if Insecure is true.
	TLSVerify = true
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
	defer func() { _ = sf.Close() }()
	df, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer func() { _ = df.Close() }()
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
	_ = os.RemoveAll(dest)
	_ = os.MkdirAll(dest, 0755)

	entries, err := ioutil.ReadDir(source)
	if err != nil {
		return fromError(err)
	}

	for _, entry := range entries {
		sfp := filepath.Join(source, entry.Name())
		dfp := filepath.Join(dest, entry.Name())
		if err := copyFile(sfp, dfp); err != nil {
			return err
		}
	}
	return nil
}

// requestURL builds the URL for the given path. The second parameter "token"
// tells this function whether it should include the authorization token in the
// query. It returns an error if this library is set to refuse insecure
// connections and a bare HTTP request is attempted.
func requestURL(path string, token bool) (string, error) {
	u, _ := url.Parse(config.Server)

	if !Insecure && u.Scheme != "https" {
		return "", errors.New("attempted to reach a server that is not using HTTPS")
	}

	u.Path = path
	if token {
		v := url.Values{}
		v.Set("token", config.Token)
		u.RawQuery = v.Encode()
	}
	return u.String(), nil
}

// safeResponse performs an HTTP request as expected by a "todo" server, while
// taking into account the TLSVerify and Insecure flags. The "method" parameter
// corresponds to an HTTP method (e.g. "GET") and the "url" parameter
// corresponds to just the path for the URL (e.g. "/topics"). Some HTTP
// requests might want to send data through the body of the request. In this
// case the "body" parameter should be used. The "token" parameter tells this
// function whether the authorization token should be sent or not with the
// request.
func safeResponse(method, url string, body io.Reader, token bool) (*http.Response, error) {
	// Setup the client for the HTTP request.
	client := http.Client{
		Timeout: requestTimeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: !TLSVerify || Insecure},
		},
	}

	str, err := requestURL(url, token)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, str, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	return client.Do(req)
}

// getResponse calls safeResponse assuming that a token is required, and then
// it polishes any given error so it can be shown to the user directly.
func getResponse(method, url string, body io.Reader) (*http.Response, error) {
	res, err := safeResponse(method, url, body, true)
	if err == nil {
		return res, nil
	}

	// Check specifically for a timeout.
	if err, ok := err.(net.Error); ok && err.Timeout() {
		return nil, NewError("timed out! Try it again in another time")
	}

	errMsg := err.Error()

	// Known limitation in Go 1.5: the *http.httpError type does not implement
	// the net.Error type, so the previous if fails in Go 1.5 but not
	// afterward.
	if strings.Contains(errMsg, "Client.Timeout") {
		return nil, NewError("timed out! Try it again in another time")
	}

	// Beautify the given error: only return the actual message.
	re, _ := regexp.Compile(`:\s+(.+)$`)
	if e := re.FindSubmatch([]byte(errMsg)); len(e) == 2 {
		return nil, NewError(string(e[1]))
	}
	return nil, fromError(err)
}
