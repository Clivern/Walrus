// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

// Migration struct
type Migration struct {
	ID    int       `json:"id"`
	Flag  string    `json:"file"`
	RunAt time.Time `json:"run_at"`
}
