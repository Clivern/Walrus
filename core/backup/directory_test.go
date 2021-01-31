// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// +build unit

package backup

import (
	"fmt"
	"testing"

	"github.com/clivern/walrus/core/util"
	"github.com/clivern/walrus/pkg"

	"github.com/franela/goblin"
)

// TestBackupRestoreDirectory
func TestBackupRestoreDirectory(t *testing.T) {
	baseDir := pkg.GetBaseDir("cache")
	pkg.LoadConfigs(fmt.Sprintf("%s/config.dist.yml", baseDir))

	g := goblin.Goblin(t)

	mgr := NewManager(nil)
	directory := fmt.Sprintf("%s/%s", baseDir, "cache/")
	target := fmt.Sprintf("%s/cache/%s", baseDir, "app.tar.gz")

	g.Describe("#TestBackupDirectory", func() {

		g.Before(func() {
			util.EnsureDir(fmt.Sprintf("%s/test/data1", directory), 0755)
			util.EnsureDir(fmt.Sprintf("%s/test/data2", directory), 0755)
			util.EnsureDir(fmt.Sprintf("%s/test/data3", directory), 0755)

			util.StoreFile(fmt.Sprintf("%s/test/data1/data.txt", directory), "data1.data\n")
			util.StoreFile(fmt.Sprintf("%s/test/data2/data.txt", directory), "data2.data\n")
			util.StoreFile(fmt.Sprintf("%s/test/data3/data.txt", directory), "data3.data\n")
			util.StoreFile(fmt.Sprintf("%s/test/data.txt", directory), "data.data\n")
		})

		g.It("It should backup", func() {
			err := mgr.Backup(
				fmt.Sprintf("%s/test/", directory),
				target,
			)

			g.Assert(err).Equal(nil)
		})

		g.It("It should restore", func() {
			util.DeleteFile(fmt.Sprintf("%s/test/data1/data.txt", directory))
			util.DeleteFile(fmt.Sprintf("%s/test/data2/data.txt", directory))
			util.DeleteFile(fmt.Sprintf("%s/test/data3/data.txt", directory))
			util.DeleteFile(fmt.Sprintf("%s/test/data.txt", directory))

			err := mgr.Restore(
				target,
				directory,
			)

			g.Assert(err).Equal(nil)

			g.Assert(util.FileExists(fmt.Sprintf("%s/test/data1/data.txt", directory))).Equal(true)
			g.Assert(util.FileExists(fmt.Sprintf("%s/test/data2/data.txt", directory))).Equal(true)
			g.Assert(util.FileExists(fmt.Sprintf("%s/test/data3/data.txt", directory))).Equal(true)
			g.Assert(util.FileExists(fmt.Sprintf("%s/test/data.txt", directory))).Equal(true)

		})

		g.After(func() {
			util.DeleteFile(target)
			util.DeleteDir(fmt.Sprintf("%s/test/", directory))
		})
	})
}
