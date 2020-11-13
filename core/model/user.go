// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// User struct
type User struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	Role      string     `json:"role"`
	LastLogin *time.Time `json:"last_login"`
}

// Users struct
type Users struct {
	Users []User `json:"users"`
}

// CreateUser creates a new user
func (db *Database) CreateUser(user *User) *User {
	db.Connection.Create(user)

	return user
}

// GetUsers gets users
func (db *Database) GetUsers() []User {
	users := []User{}

	db.Connection.Select("*").Find(&users)

	return users
}

// GetUserByEmail gets a user by email
func (db *Database) GetUserByEmail(email string) User {
	user := User{}

	db.Connection.Where("email = ?", email).First(&user)

	return user
}
