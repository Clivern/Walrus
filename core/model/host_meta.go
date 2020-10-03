// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
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

// LoadFromJSON update object from json
func (h *HostMeta) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &h)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (h *HostMeta) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&h)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (h *HostMetas) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &h)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (h *HostMetas) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&h)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
