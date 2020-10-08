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
