// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// HostMeta struct
type HostMeta struct {
	ID        int       `json:"id"`
	HostID    int       `json:"host_id"`
	Key       string    `json:"key"`
	Value     string    `json:"value"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// HostMetas struct
type HostMetas struct {
	HostMetas []HostMeta `json:"hosts"`
}
