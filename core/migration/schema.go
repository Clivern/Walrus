// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package migration

import (
	"time"

	"github.com/jinzhu/gorm"
)

// Option struct
type Option struct {
	gorm.Model

	Key   string `json:"key"`
	Value string `json:"value"`
}

// User struct
type User struct {
	gorm.Model

	Role      string    `json:"role"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	LastLogin time.Time `json:"last_login"`
}

// Job struct
type Job struct {
	gorm.Model

	UUID    string    `json:"uuid"`
	Payload string    `json:"payload"`
	Status  string    `json:"status"`
	Type    string    `json:"type"`
	Result  string    `json:"result"`
	Retry   int       `json:"retry"`
	Parent  int       `json:"parent"`
	HostID  int       `json:"host_id"`
	RunAt   time.Time `json:"run_at"`
}

// Host struct
type Host struct {
	gorm.Model

	Name            string    `json:"name"`
	UUID            string    `json:"uuid"`
	RetentionPolicy string    `json:"retention_policy"`
	StorageID       string    `json:"storage_id"`
	Status          string    `json:"status"`
	LastCheck       time.Time `json:"last_check"`
}

// HostMeta struct
type HostMeta struct {
	gorm.Model

	HostID int    `json:"host_id"`
	Key    string `json:"key"`
	Value  string `json:"value"`
}
