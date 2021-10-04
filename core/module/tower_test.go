// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//go:build integration
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

// TestTowerMethods test cases
func TestTowerMethods(t *testing.T) {
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

	tower := NewTower(db)

	// Cleanup
	db.Delete(viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))))

	time.Sleep(3 * time.Second)

	g.Describe("#Elect", func() {
		g.It("It should satisfy test cases", func() {
			err := tower.Elect(int64(6))

			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#IsLeader", func() {
		g.It("It should satisfy test cases", func() {
			ok, err := tower.IsLeader()

			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#HasLeader", func() {
		g.It("It should satisfy test cases", func() {
			ok, err := tower.HasLeader()

			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#Alive", func() {
		g.It("It should satisfy test cases", func() {
			err := tower.Alive(int64(5))

			g.Assert(err).Equal(nil)
		})
	})
}
