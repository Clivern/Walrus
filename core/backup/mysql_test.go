// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

// +build unit

package backup

import (
	"testing"

	"github.com/franela/goblin"
)

// TestMySQLType
func TestMySQLType(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("#TestDumpOptions", func() {

		g.It("It should return expected options", func() {
			mysql := &MySQL{
				Host:         "127.0.0.1",
				Port:         "3306",
				Username:     "root",
				Password:     "root",
				AllDatabases: true,
				Options:      "--single-transaction,--quick,--lock-tables=false",
				OutputFile:   "/tmp/result.sql",
			}

			g.Assert(mysql.DumpOptions()).Equal("--host 127.0.0.1 --port 3306 -u root -proot --result-file=/tmp/result.sql --all-databases --single-transaction --quick --lock-tables=false")
		})

		g.It("It should return expected options", func() {
			mysql := &MySQL{
				Host:         "127.0.0.1",
				Port:         "3306",
				Username:     "root",
				Password:     "root",
				AllDatabases: false,
				Database:     "walrus",
				Table:        "",
				Options:      "--single-transaction,--quick,--lock-tables=false",
				OutputFile:   "/tmp/result.sql",
			}

			g.Assert(mysql.DumpOptions()).Equal("--host 127.0.0.1 --port 3306 -u root -proot --result-file=/tmp/result.sql walrus --single-transaction --quick --lock-tables=false")
		})

		g.It("It should return expected options", func() {
			mysql := &MySQL{
				Host:         "127.0.0.1",
				Port:         "3306",
				Username:     "root",
				Password:     "root",
				AllDatabases: false,
				Database:     "walrus",
				Table:        "job",
				Options:      "--single-transaction,--quick,--lock-tables=false",
				OutputFile:   "/tmp/result.sql",
			}

			g.Assert(mysql.DumpOptions()).Equal("--host 127.0.0.1 --port 3306 -u root -proot --result-file=/tmp/result.sql walrus job --single-transaction --quick --lock-tables=false")
		})
	})
}
