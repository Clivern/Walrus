// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"testing"

	"github.com/clivern/walrus/pkg"
)

// TestDsnToString test cases
func TestDsnToString(t *testing.T) {
	t.Run("TestDsnToStringForMySQL", func(t *testing.T) {
		dsn := DSN{
			Driver:   "mysql",
			Username: "root",
			Password: "root",
			Hostname: "127.0.0.1",
			Port:     3306,
			Name:     "walrus",
		}
		pkg.Expect(t, "root:root@tcp(127.0.0.1:3306)/walrus?charset=utf8&parseTime=True", dsn.ToString())
	})

	t.Run("TestDsnToStringForSQLLite", func(t *testing.T) {
		dsn := DSN{
			Driver: "sqlite3",
			Name:   "/path/to/walrus.db",
		}
		pkg.Expect(t, "/path/to/walrus.db", dsn.ToString())
	})
}
