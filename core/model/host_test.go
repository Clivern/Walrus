// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// +build integration

package model

import (
	"fmt"
	"testing"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/pkg"

	"github.com/franela/goblin"
	"github.com/spf13/viper"
)

// TestHostMethods test cases
func TestHostMethods(t *testing.T) {
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

	host := NewHostStore(db)

	// Cleanup
	db.Delete(viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))))

	time.Sleep(3 * time.Second)

	g.Describe("#CreateHost", func() {
		g.It("It should satisfy test cases", func() {
			err := host.CreateHost(HostData{
				Hostname: "localhost",
			})

			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#UpdateHost", func() {
		g.It("It should satisfy test cases", func() {
			err := host.UpdateHost(HostData{
				Hostname: "localhost",
			})

			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#GetHost", func() {
		g.It("It should satisfy test cases", func() {
			value, err := host.GetHost("localhost")

			g.Assert(value.Hostname).Equal("localhost")
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#GetHosts", func() {
		g.It("It should satisfy test cases", func() {
			value, err := host.GetHosts()

			g.Assert(value[0].Hostname).Equal("localhost")
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#DeleteHost", func() {
		g.It("It should satisfy test cases", func() {
			ok, err := host.DeleteHost("localhost")

			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)
		})
	})
}
