// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// +build integration

package module

import (
	"fmt"
	"testing"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/pkg"

	"github.com/franela/goblin"
	"github.com/spf13/viper"
)

// TestStatsMethods test cases
func TestStatsMethods(t *testing.T) {
	// Skip if -short flag exist
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	baseDir := pkg.GetBaseDir("cache")
	pkg.LoadConfigs(fmt.Sprintf("%s/config.dist.yml", baseDir))

	g := goblin.Goblin(t)

	db := driver.NewEtcdDriver()
	db.Connect()

	defer db.Close()

	stats := NewStats(db)

	// Cleanup
	db.Delete(viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))))

	time.Sleep(3 * time.Second)

	g.Describe("#GetTotalTowers", func() {
		g.It("It should satisfy test cases", func() {
			count, err := stats.GetTotalTowers()

			g.Assert(count).Equal(0)
			g.Assert(err).Equal(nil)

			db.Put(fmt.Sprintf(
				"%s/tower/a/data",
				viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
			), "#")

			db.Put(fmt.Sprintf(
				"%s/tower/b/data",
				viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
			), "#")

			count, err = stats.GetTotalTowers()

			g.Assert(count).Equal(2)
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#GetTotalHosts", func() {
		g.It("It should satisfy test cases", func() {
			count, err := stats.GetTotalHosts()

			g.Assert(count).Equal(0)
			g.Assert(err).Equal(nil)

			db.Put(fmt.Sprintf(
				"%s/host/a/agent/a/data",
				viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
			), "#")

			db.Put(fmt.Sprintf(
				"%s/host/a/agent/b/data",
				viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
			), "#")

			db.Put(fmt.Sprintf(
				"%s/host/b/agent/b/data",
				viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
			), "#")

			count, err = stats.GetTotalHosts()

			g.Assert(count).Equal(2)
			g.Assert(err).Equal(nil)
		})
	})
}
