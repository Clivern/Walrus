// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"os"

	"github.com/BurntSushi/toml"
)

// Configs type
type Configs struct {
	General General `toml:"General"`
}

// General type
type General struct {
	Key string `toml:"key"`
}

// NewConfigs creates an instance of Configs
func NewConfigs() *Configs {
	return &Configs{
		General: General{
			Key: "value",
		},
	}
}

// Decode decodes from file to struct
func (g *Configs) Decode(path string) error {
	if _, err := toml.DecodeFile(path, &g); err != nil {
		return err
	}

	return nil
}

// Encode encodes struct and store on file
func (g *Configs) Encode(path string) error {
	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	err = toml.NewEncoder(f).Encode(g)

	if err != nil {
		return err
	}

	return nil
}
