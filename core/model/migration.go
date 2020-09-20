// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"
)

// Migration struct
type Migration struct {
	ID    int       `json:"id"`
	Flag  string    `json:"file"`
	RunAt time.Time `json:"run_at"`
}

// LoadFromJSON update object from json
func (m *Migration) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (m *Migration) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
