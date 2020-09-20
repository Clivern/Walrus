// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
)

// Message struct
type Message struct {
	UUID string `json:"uuid"`
	Job  int    `json:"job"`
}

// LoadFromJSON update object from json
func (m *Message) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &m)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (m *Message) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&m)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
