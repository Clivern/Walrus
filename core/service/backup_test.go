// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/clivern/walrus/pkg"
)

// TestBackupRestoreDirectory test cases
func TestBackupRestoreDirectory(t *testing.T) {
	dir, _ := os.Getwd()
	cache := fmt.Sprintf("%s/%s", dir, "cache")

	for {
		if DirExists(cache) {
			break
		}
		dir = filepath.Dir(dir)
		cache = fmt.Sprintf("%s/%s", dir, "cache")
	}

	backupClient := &Backup{}
	directory := fmt.Sprintf("%s/%s", dir, "cache/")
	target := fmt.Sprintf("%s/%s", cache, "app.tar.gz")

	t.Run("CreateTestDir", func(t *testing.T) {
		EnsureDir(fmt.Sprintf("%s/test/data1", directory), 0755)
		EnsureDir(fmt.Sprintf("%s/test/data2", directory), 0755)
		EnsureDir(fmt.Sprintf("%s/test/data3", directory), 0755)

		StoreFile(fmt.Sprintf("%s/test/data1/data.txt", directory), "data1.data\n")
		StoreFile(fmt.Sprintf("%s/test/data2/data.txt", directory), "data2.data\n")
		StoreFile(fmt.Sprintf("%s/test/data3/data.txt", directory), "data3.data\n")
		StoreFile(fmt.Sprintf("%s/test/data.txt", directory), "data.data\n")
	})

	t.Run("TestBackupDirectory", func(t *testing.T) {
		ok, err := backupClient.BackupDirectory(
			fmt.Sprintf("%s/test/", directory),
			target,
		)
		pkg.Expect(t, ok, true)
		pkg.Expect(t, err, nil)
	})

	t.Run("ChangeTestDir", func(t *testing.T) {
		DeleteFile(fmt.Sprintf("%s/test/data1/data.txt", directory))
		DeleteFile(fmt.Sprintf("%s/test/data2/data.txt", directory))
		DeleteFile(fmt.Sprintf("%s/test/data3/data.txt", directory))
		DeleteFile(fmt.Sprintf("%s/test/data.txt", directory))
	})

	t.Run("TestRestoreDirectory", func(t *testing.T) {
		ok, err := backupClient.RestoreDirectory(
			target,
			directory,
		)
		pkg.Expect(t, ok, true)
		pkg.Expect(t, err, nil)
	})

	t.Run("ValidateTestDir", func(t *testing.T) {
		pkg.Expect(t, FileExists(fmt.Sprintf("%s/test/data1/data.txt", directory)), true)
		pkg.Expect(t, FileExists(fmt.Sprintf("%s/test/data2/data.txt", directory)), true)
		pkg.Expect(t, FileExists(fmt.Sprintf("%s/test/data3/data.txt", directory)), true)
		pkg.Expect(t, FileExists(fmt.Sprintf("%s/test/data.txt", directory)), true)
	})

	t.Run("Cleanup", func(t *testing.T) {
		DeleteFile(target)
		DeleteDir(fmt.Sprintf("%s/test/", directory))
	})
}
