// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"github.com/clivern/walrus/core/driver"
)

// Backup type
type Backup struct {
	db driver.Database
}

// NewBackupStore creates a new instance
func NewBackupStore(db driver.Database) *Backup {
	result := new(Backup)
	result.db = db

	return result
}
