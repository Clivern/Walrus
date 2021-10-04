// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

//go:build unit
// +build unit

package util

import (
	"fmt"
	"testing"

	"github.com/clivern/walrus/pkg"

	"github.com/franela/goblin"
)

// TestHelpers
func TestHelpers(t *testing.T) {
	baseDir := pkg.GetBaseDir("cache")
	pkg.LoadConfigs(fmt.Sprintf("%s/config.dist.yml", baseDir))

	g := goblin.Goblin(t)

	g.Describe("#TestInArray", func() {
		g.It("It should satisfy test cases", func() {
			g.Assert(InArray("A", []string{"A", "B", "C", "D"})).Equal(true)
			g.Assert(InArray("B", []string{"A", "B", "C", "D"})).Equal(true)
			g.Assert(InArray("H", []string{"A", "B", "C", "D"})).Equal(false)
			g.Assert(InArray(1, []int{2, 3, 1})).Equal(true)
			g.Assert(InArray(9, []int{2, 3, 1})).Equal(false)
		})
	})

	g.Describe("#TestEncryption", func() {
		g.It("It should satisfy test cases", func() {
			ciphertext, err := Encrypt([]byte("Hello World"), "password")
			plaintext, err := Decrypt(ciphertext, "password")

			g.Assert(string(plaintext)).Equal("Hello World")
			g.Assert(err).Equal(nil)

			ciphertext = []byte("Hello World")
			plaintext, err = Decrypt(ciphertext, "password")

			g.Assert(string(plaintext)).Equal("")
			g.Assert(err.Error()).Equal("Invalid encrypted text")
		})
	})

	g.Describe("#TestHashing", func() {
		g.It("It should satisfy test cases", func() {
			hash, err := HashPassword("password")
			g.Assert(true).Equal(CheckPasswordHash("password", hash))
			g.Assert(false).Equal(CheckPasswordHash("PasSword", hash))
			g.Assert(err).Equal(nil)
		})
	})
}
