// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"time"
)

var (
	// JobPending pending job type
	JobPending = "pending"

	// JobFailed failed job type
	JobFailed = "failed"

	// JobSuccess success job type
	JobSuccess = "success"

	// JobOnHold on hold job type
	JobOnHold = "on_hold"
)

// Job struct
type Job struct {
	ID        int        `json:"id"`
	UUID      string     `json:"uuid"`
	Payload   string     `json:"payload"`
	Status    string     `json:"status"`
	Type      string     `json:"type"`
	Result    string     `json:"result"`
	Retry     int        `json:"retry"`
	Parent    int        `json:"parent"`
	RunAt     *time.Time `json:"run_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Jobs struct
type Jobs struct {
	Jobs []Job `json:"jobs"`
}
