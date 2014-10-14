// Copyright (C) 2014 Miquel Sabaté Solà <mikisabate@gmail.com>
// This file is licensed under the MIT license.
// See the LICENSE file.

package lib

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	defaultEditor = "vi"
	dirName       = ".td"
	fileName      = ".config.json"
	tmpDir        = "tmp"
	oldDir        = "old"
	newDir        = "new"
)

// TODO: establish a locking mechanism

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

func editor() string {
	value := os.Getenv("EDITOR")
	if value == "" {
		return defaultEditor
	}
	return value
}

// TODO
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

// TODO
func copyDir(source string, dest string) error {
	os.RemoveAll(dest)
	os.MkdirAll(dest, 0755)

	entries, err := ioutil.ReadDir(source)
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

func save(topics []Topic) error {
	// First of all, reset the temporary directory.
	dir := filepath.Join(home(), dirName, tmpDir)
	os.RemoveAll(dir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fromError(err)
	}

	// Save all the topics to this temporary directory.
	for _, t := range topics {
		path := filepath.Join(dir, t.Name+".md")
		if err := write(&t, path); err != nil {
			return err
		}
	}

	// And update the old and new directories
	adir := filepath.Join(home(), dirName, oldDir)
	if err := copyDir(dir, adir); err != nil {
		return fromError(err)
	}
	adir = filepath.Join(home(), dirName, newDir)
	if err := copyDir(dir, adir); err != nil {
		return fromError(err)
	}

	return nil
}

func write(topic *Topic, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fromError(err)
	}
	defer f.Close()
	if _, err := f.WriteString(topic.Contents); err != nil {
		return fromError(err)
	}
	return nil
}
