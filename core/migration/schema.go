// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package migration

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// Option struct
type Option struct {
	gorm.Model

	Key   string `json:"key"`
	Value string `json:"value"`
}

// Job struct
type Job struct {
	gorm.Model

	UUID    string    `json:"uuid"`
	Payload string    `json:"payload"`
	Status  string    `json:"status"`
	Type    string    `json:"type"`
	Result  string    `json:"result"`
	Retry   int       `json:"retry"`
	Parent  int       `json:"parent"`
	RunAt   time.Time `json:"run_at"`
}

// Agent struct
type Agent struct {
	gorm.Model

	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	URL       string    `json:"url"`
	LastCheck time.Time `json:"last_check"`
}

// Host struct
type Host struct {
	gorm.Model

	UUID      string    `json:"uuid"`
	Configs   string    `json:"configs"`
	Status    string    `json:"status"`
	Type      string    `json:"type"`
	DestroyAt time.Time `json:"destroy_at"`
}

// LoadFromJSON update object from json
func (o *Option) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &o)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (o *Option) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&o)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// LoadFromJSON update object from json
func (j *Job) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &j)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (j *Job) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&j)
	if err != nil {
		return "", err
	}
	return string(data), nil
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
func (h *Host) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &h)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (h *Host) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&h)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
