// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
	"time"
)

// Agent struct
type Agent struct {
	ID        int        `json:"id"`
	UUID      string     `json:"uuid"`
	Name      string     `json:"name"`
	Status    string     `json:"status"`
	URL       string     `json:"url"`
	LastCheck *time.Time `json:"last_check"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Agents struct
type Agents struct {
	Agents []Agent `json:"agents"`
}

// LoadFromJSON update object from json
func (a *Agent) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &a)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (a *Agent) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (a *Agents) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &a)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (a *Agents) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
