// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"
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

// LoadFromJSON update object from json
func (s *Host) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &s)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (s *Host) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (s *Hosts) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &s)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (s *Hosts) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&s)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
