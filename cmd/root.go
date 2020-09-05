// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/clivern/walrus/core/util"

	"github.com/spf13/cobra"
)

// Verbose var
var Verbose bool

// Config var
var Config string

// ConfigFilePath var
const ConfigFilePath = "walrus.toml"

var rootCmd = &cobra.Command{
	Use: "walrus",
	Short: `🔥 A fast and beautiful terminal recorder

If you have any suggestions, bug reports, or annoyances please report
them to our issue tracker at <https://github.com/clivern/walrus/issues>`,
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().StringVarP(
		&Config,
		"config",
		"c",
		fmt.Sprintf("%s%s", util.EnsureTrailingSlash(os.Getenv("HOME")), ConfigFilePath),
		"config file",
	)
}

// Execute runs cmd tool
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
