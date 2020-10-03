// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"
)

// User struct
type User struct {
	ID        int        `json:"id"`
	Email     string     `json:"email"`
	Password  string     `json:"password"`
	LastLogin *time.Time `json:"last_login"`
}

// Users struct
type Users struct {
	Users []User `json:"users"`
}

// LoadFromJSON update object from json
func (u *User) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &u)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (u *User) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&u)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (u *Users) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &u)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (u *Users) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&u)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
