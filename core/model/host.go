// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"github.com/clivern/walrus/core/migration"
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

// CreateHost creates a new host
func (db *Database) CreateHost(host *Host) *Host {
	db.Connection.Create(host)
	return host
}

// HostExistByID check if host exists
func (db *Database) HostExistByID(id int) bool {
	host := Host{}

	db.Connection.Where("id = ?", id).First(&host)

	return host.ID > 0
}

// GetHostByID gets a host by id
func (db *Database) GetHostByID(id int) Host {
	host := Host{}

	db.Connection.Where("id = ?", id).First(&host)

	return host
}

// GetHosts gets services
func (db *Database) GetHosts() []Host {
	services := []Host{}

	db.Connection.Select("*").Find(&services)

	return services
}

// HostExistByUUID check if host exists
func (db *Database) HostExistByUUID(uuid string) bool {
	host := Host{}

	db.Connection.Where("uuid = ?", uuid).First(&host)

	return host.ID > 0
}

// GetHostByUUID gets a host by uuid
func (db *Database) GetHostByUUID(uuid string) Host {
	host := Host{}

	db.Connection.Where("uuid = ?", uuid).First(&host)

	return host
}

// DeleteHostByID deletes a host by id
func (db *Database) DeleteHostByID(id int) {
	db.Connection.Unscoped().Where("id=?", id).Delete(&migration.Host{})
	db.Connection.Unscoped().Where("host_id=?", id).Delete(&migration.Job{})
	db.Connection.Unscoped().Where("host_id=?", id).Delete(&migration.HostMeta{})
}

// DeleteHostByUUID deletes a host by uuid
func (db *Database) DeleteHostByUUID(uuid string) {
	host := db.GetHostByUUID(uuid)

	if host.ID > 0 {
		db.DeleteHostByID(host.ID)
	}
}

// UpdateHostByID updates a host by ID
func (db *Database) UpdateHostByID(host *Host) {
	db.Connection.Save(&host)
}
