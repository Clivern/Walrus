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

// TestAgentMethods test cases
func TestAgentMethods(t *testing.T) {
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

	agent := NewAgentStore(db)

	// Cleanup
	db.Delete(viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))))

	time.Sleep(3 * time.Second)

	g.Describe("#CreateAgent", func() {
		g.It("It should satisfy test cases", func() {
			err := agent.CreateAgent(AgentData{
				ID:        "41aae480-e338-4b5a-86bb-1f5c8b15ab7b",
				URL:       "http://localhost:8080",
				Hostname:  "x-hostname",
				APIKey:    "x-api-key",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			})

			g.Assert(err).Equal(nil)

			err = agent.CreateAgent(AgentData{
				ID:        "70f8c7f6-ca1e-4f2b-9f74-b99f8789b023",
				URL:       "http://localhost:8081",
				Hostname:  "y-hostname",
				APIKey:    "y-api-key",
				CreatedAt: time.Now().Unix(),
				UpdatedAt: time.Now().Unix(),
			})

			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#UpdateAgent", func() {
		g.It("It should satisfy test cases", func() {
			err := agent.UpdateAgent(AgentData{
				ID:        "70f8c7f6-ca1e-4f2b-9f74-b99f8789b023",
				URL:       "http://localhost:8089",
				Hostname:  "y-hostname",
				APIKey:    "z-api-key",
				CreatedAt: time.Now().Unix(),
			})

			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#GetAgent", func() {
		g.It("It should satisfy test cases", func() {
			value, err := agent.GetAgent("y-hostname", "70f8c7f6-ca1e-4f2b-9f74-b99f8789b023")

			g.Assert(value.URL).Equal("http://localhost:8089")
			g.Assert(err).Equal(nil)

			value, err = agent.GetAgent("x-hostname", "41aae480-e338-4b5a-86bb-1f5c8b15ab7b")

			g.Assert(value.APIKey).Equal("x-api-key")
			g.Assert(err).Equal(nil)
		})
	})

	g.Describe("#DeleteAgent", func() {
		g.It("It should satisfy test cases", func() {
			ok, err := agent.DeleteAgent("y-hostname", "70f8c7f6-ca1e-4f2b-9f74-b99f8789b023")

			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)

			ok, err = agent.DeleteAgent("x-hostname", "41aae480-e338-4b5a-86bb-1f5c8b15ab7b")

			g.Assert(ok).Equal(true)
			g.Assert(err).Equal(nil)
		})
	})
}
