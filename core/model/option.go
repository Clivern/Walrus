// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"github.com/clivern/walrus/core/migration"
)

// Option struct
type Option struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Options struct
type Options struct {
	Options []Option `json:"options"`
}

// CreateOption creates a new option
func (db *Database) CreateOption(option *Option) *Option {
	db.Connection.Create(option)
	return option
}

// OptionExistByID check if option exists
func (db *Database) OptionExistByID(id int) bool {
	option := Option{}

	db.Connection.Where("id = ?", id).First(&option)

	return option.ID > 0
}

// OptionExistByKey check if an option exists
func (db *Database) OptionExistByKey(key string) bool {
	option := Option{}

	db.Connection.Where("key = ?", key).First(&option)

	return option.ID > 0
}

// GetOptionByID gets an option by id
func (db *Database) GetOptionByID(id int) Option {
	option := Option{}

	db.Connection.Where("id = ?", id).First(&option)

	return option
}

// GetOptionByKey gets an option by key
func (db *Database) GetOptionByKey(key string) Option {
	option := Option{}

	db.Connection.Where("key = ?", key).First(&option)

	return option
}

// GetOptions gets options
func (db *Database) GetOptions() []Option {
	options := []Option{}

	db.Connection.Select("*").Find(&options)

	return options
}

// DeleteOptionByID deletes an option by id
func (db *Database) DeleteOptionByID(id int) {
	db.Connection.Unscoped().Where("id=?", id).Delete(&migration.Option{})
}

// DeleteOptionByKey deletes an option by key
func (db *Database) DeleteOptionByKey(key string) {
	db.Connection.Unscoped().Where("key=?", key).Delete(&migration.Option{})
}

// UpdateOptionByID updates an option by ID
func (db *Database) UpdateOptionByID(option *Option) {
	db.Connection.Save(&option)
}
