// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"encoding/json"
)

// Option struct
type Option struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// Options struct
type Options struct {
	Options []Option `json:"options"`
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
func (o *Options) LoadFromJSON(data []byte) (bool, error) {
	err := json.Unmarshal(data, &o)
	if err != nil {
		return false, err
	}
	return true, nil
}

// ConvertToJSON convert object to json
func (o *Options) ConvertToJSON() (string, error) {
	data, err := json.Marshal(&o)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
