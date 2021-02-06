// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package backup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/clivern/walrus/core/util"
)

// BackupDirectory creates a backup of any directory
// directory must be absolute path like /etc/app and archive same /etc/app.tar.gz"
func (m *Manager) BackupDirectory(directory, archive string) error {
	if !util.DirExists(directory) && !util.FileExists(directory) {
		return fmt.Errorf("Unable to find directory or file %s", directory)
	}
	// Get the full executable path for the editor.
	executable, err := exec.LookPath("tar")

	if err != nil {
		return err
	}

	command := strings.Split(
		fmt.Sprintf(`%s -czf %s %s`, executable, archive, directory),
		" ",
	)

	cmd := exec.Command(command[0], command[1:]...)

	cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RestoreDirectory uncompress .tar.gz archive to the root path /
func (m *Manager) RestoreDirectory(archive, destination string) error {
	if !util.FileExists(archive) {
		return fmt.Errorf("Unable to find archive %s", archive)
	}

	// Get the full executable path for the editor.
	executable, err := exec.LookPath("tar")

	if err != nil {
		return err
	}

	command := strings.Split(
		fmt.Sprintf(`%s -xzf %s -C %s`, executable, archive, destination),
		" ",
	)

	cmd := exec.Command(command[0], command[1:]...)

	cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	return cmd.Run()
}
