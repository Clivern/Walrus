// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// PathExists reports whether the path exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// FileExists reports whether the named file exists
func FileExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}

	return false
}

// DirExists reports whether the dir exists
func DirExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}

	return false
}

// EnsureDir ensures that directory exists
func EnsureDir(dirName string, mode int) (bool, error) {
	err := os.MkdirAll(dirName, os.FileMode(mode))

	if err == nil || os.IsExist(err) {
		return true, nil
	}

	return false, err
}

// DeleteFile deletes a file
func DeleteFile(path string) bool {
	if err := os.Remove(path); err == nil {
		return true
	}

	return false
}

// ValidRelPath checks for path traversal and correct forward slashes
func ValidRelPath(path string) bool {
	if path == "" || strings.Contains(path, `\`) || strings.HasPrefix(path, "/") || strings.Contains(path, "../") {
		return false
	}

	return true
}

// EnsureTrailingSlash ensure there is a trailing slash
func EnsureTrailingSlash(dir string) string {
	return fmt.Sprintf(
		"%s%s",
		strings.TrimRight(dir, string(os.PathSeparator)),
		string(os.PathSeparator),
	)
}

// RemoveTrailingSlash removes any trailing slash
func RemoveTrailingSlash(dir string) string {
	return strings.TrimRight(dir, string(os.PathSeparator))
}

// RemoveStartingSlash removes any starting slash
func RemoveStartingSlash(dir string) string {
	return strings.TrimLeft(dir, string(os.PathSeparator))
}

// ClearDir removes all files and sub dirs
func ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteDir deletes a dir
func DeleteDir(dir string) bool {
	if err := os.RemoveAll(dir); err == nil {
		return true
	}

	return false
}

// StoreFile stores a file content
func StoreFile(path, content string) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, 0775)

	if err != nil {
		return err
	}

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)

	return err
}
