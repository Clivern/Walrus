// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// User type
type User struct {
	db driver.Database
}

// UserData type
type UserData struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	PasswordHash string `json:"passwordHash"`
	APIKey       string `json:"apiKey"`
	CreatedAt    int64  `json:"createdAt"`
	UpdatedAt    int64  `json:"updatedAt"`
}

// NewUserStore creates a new instance
func NewUserStore(db driver.Database) *User {
	result := new(User)
	result.db = db

	return result
}

// CreateUser creates a user
func (u *User) CreateUser(user UserData) error {
	// Generate password hash
	var err error

	if user.Password == "" && user.PasswordHash == "" {
		return fmt.Errorf("Error! both password and password hash are missing")
	}

	if user.Password != "" {
		user.PasswordHash, err = util.HashPassword(user.Password)

		if err != nil {
			return err
		}

		user.Password = ""
	}

	user.CreatedAt = time.Now().Unix()
	user.UpdatedAt = time.Now().Unix()

	result, err := util.ConvertToJSON(user)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"user_id":    user.ID,
		"user_email": user.Email,
	}).Debug("Create a user")

	// store user data
	err = u.db.Put(fmt.Sprintf(
		"%s/user/%s/u-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		user.Email,
	), result)

	if err != nil {
		return err
	}

	return nil
}

// UpdateUserByEmail updates a user by email
func (u *User) UpdateUserByEmail(user UserData) error {
	var err error

	if user.Password == "" && user.PasswordHash == "" {
		return fmt.Errorf("Error! both password and password hash are missing")
	}

	if user.Password != "" {
		user.PasswordHash, err = util.HashPassword(user.Password)

		if err != nil {
			return err
		}

		user.Password = ""
	}

	user.UpdatedAt = time.Now().Unix()

	result, err := util.ConvertToJSON(user)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"user_id":    user.ID,
		"user_email": user.Email,
	}).Debug("Update user")

	// store user data
	err = u.db.Put(fmt.Sprintf(
		"%s/user/%s/u-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		user.Email,
	), result)

	if err != nil {
		return err
	}

	return nil
}

// GetUserByEmail gets a user by email
func (u *User) GetUserByEmail(email string) (*UserData, error) {
	user := &UserData{}

	log.WithFields(log.Fields{
		"user_email": email,
	}).Debug("Get a user")

	data, err := u.db.Get(fmt.Sprintf(
		"%s/user/%s/u-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		email,
	))

	if err != nil {
		return user, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/u-data") {
			err = util.LoadFromJSON(user, []byte(v))

			if err != nil {
				return user, err
			}

			return user, nil
		}
	}

	return user, fmt.Errorf(
		"Unable to find user with email: %s",
		email,
	)
}

// DeleteUserByEmail deletes a user by email
func (u *User) DeleteUserByEmail(email string) (bool, error) {

	log.WithFields(log.Fields{
		"user_email": email,
	}).Debug("Delete a user")

	count, err := u.db.Delete(fmt.Sprintf(
		"%s/user/%s",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		email,
	))

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// Authenticate authenticates a user
func (u *User) Authenticate(email, password string) (bool, error) {

	log.WithFields(log.Fields{
		"user_email": email,
	}).Debug("Authenticate a user")

	user, err := u.GetUserByEmail(email)

	if err != nil {
		return false, err
	}

	ok := util.CheckPasswordHash(password, user.PasswordHash)

	if !ok {
		return false, nil
	}

	return true, nil
}
