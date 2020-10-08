// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// Host struct
type Host struct {
	ID              int        `json:"id"`
	Name            string     `json:"name"`
	UUID            string     `json:"uuid"`
	RetentionPolicy string     `json:"retention_policy"`
	StorageID       string     `json:"storage_id"`
	Status          string     `json:"status"`
	LastCheck       *time.Time `json:"last_check"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// Hosts struct
type Hosts struct {
	Hosts []Host `json:"hosts"`
}
