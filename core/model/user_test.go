// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//go:build integration
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

// TestUserMethods test cases
func TestUserMethods(t *testing.T) {
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

	user := NewUserStore(db)

	// Cleanup
	db.Delete(viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))))

	time.Sleep(3 * time.Second)

	g.Describe("#CreateUser", func() {
		g.It("It should satisfy test cases", func() {
			err := user.CreateUser(UserData{
				ID:           "36e9f46d-bc1e-430f-8e5d-c6f56f0e8f0d",
				Name:         "Clivern",
				Email:        "test@clivern.com",
				Password:     "",
				PasswordHash: "",
				APIKey:       "bf67fae3-6bf2-4f2d-bb5a-67ae0fca8c97",
			})
			g.Assert(err.Error()).Equal("Error! both password and password hash are missing")

			err = user.CreateUser(UserData{
				ID:           "36e9f46d-bc1e-430f-8e5d-c6f56f0e8f0d",
				Name:         "Clivern",
				Email:        "test@clivern.com",
				Password:     "123456$!@#fw",
				PasswordHash: "",
				APIKey:       "bf67fae3-6bf2-4f2d-bb5a-67ae0fca8c97",
			})
			g.Assert(err).Equal(nil)

			err = user.CreateUser(UserData{
				ID:           "37ead05a-44b0-4b09-985d-17a22c6d3d99",
				Name:         "Joe",
				Email:        "joe@clivern.com",
				Password:     "$@gt!&*kil^yt",
				PasswordHash: "",
				APIKey:       "4beb7506-0057-4ae4-8f7e-ea2aa027b7e7",
			})
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#Authenticate", func() {
		g.It("It should satisfy test cases", func() {
			ok, err := user.Authenticate("test@clivern.com", "123456$!@#fw")
			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)

			ok, err = user.Authenticate("joe@clivern.com", "$@gt!&*kil^yt")
			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#UpdateUserByEmail", func() {
		g.It("It should satisfy test cases", func() {
			err := user.UpdateUserByEmail(UserData{
				ID:           "36e9f46d-bc1e-430f-8e5d-c6f56f0e8f0t",
				Name:         "Clivern",
				Email:        "test@clivern.com",
				Password:     "",
				PasswordHash: "",
				APIKey:       "bf67fae3-6bf2-4f2d-bb5a-67ae0fca8c98",
			})
			g.Assert(err.Error()).Equal("Error! both password and password hash are missing")

			err = user.UpdateUserByEmail(UserData{
				ID:           "36e9f46d-bc1e-430f-8e5d-c6f56f0e8f0t",
				Name:         "Clivern",
				Email:        "test@clivern.com",
				Password:     "123456$!@#fw",
				PasswordHash: "",
				APIKey:       "bf67fae3-6bf2-4f2d-bb5a-67ae0fca8c98",
			})
			g.Assert(err).Equal(nil)

			err = user.UpdateUserByEmail(UserData{
				ID:           "37ead05a-44b0-4b09-985d-17a22c6d3d93",
				Name:         "Joe",
				Email:        "joe@clivern.com",
				Password:     "$@gt!&*kil^yt",
				PasswordHash: "",
				APIKey:       "4beb7506-0057-4ae4-8f7e-ea2aa027b7e8",
			})
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#GetUserByEmail", func() {
		g.It("It should satisfy test cases", func() {
			userData, err := user.GetUserByEmail("test@clivern.com")

			g.Assert(userData.APIKey).Equal("bf67fae3-6bf2-4f2d-bb5a-67ae0fca8c98")
			g.Assert(userData.ID).Equal("36e9f46d-bc1e-430f-8e5d-c6f56f0e8f0t")
			g.Assert(err).Equal(nil)

			userData, err = user.GetUserByEmail("joe@clivern.com")

			g.Assert(userData.APIKey).Equal("4beb7506-0057-4ae4-8f7e-ea2aa027b7e8")
			g.Assert(userData.ID).Equal("37ead05a-44b0-4b09-985d-17a22c6d3d93")
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#DeleteUserByEmail", func() {
		g.It("It should satisfy test cases", func() {
			ok, err := user.DeleteUserByEmail("test@clivern.com")
			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)

			ok, err = user.DeleteUserByEmail("joe@clivern.com")
			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)
		})
	})
}
