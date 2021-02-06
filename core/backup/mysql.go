// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package backup

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// MySQL type
type MySQL struct {
	Host         string
	Port         string
	Username     string
	Password     string
	AllDatabases bool
	Database     string
	Table        string
	Options      string
	OutputFile   string
}

// DumpOptions dump backup options
func (m *MySQL) DumpOptions() string {
	if m.AllDatabases {
		return fmt.Sprintf(
			"--host %s --port %s -u %s -p%s --result-file=%s --all-databases %s",
			m.Host,
			m.Port,
			m.Username,
			m.Password,
			m.OutputFile,
			strings.Replace(m.Options, ",", " ", -1),
		)
	} else if m.Table == "" {
		return fmt.Sprintf(
			"--host %s --port %s -u %s -p%s --result-file=%s %s %s",
			m.Host,
			m.Port,
			m.Username,
			m.Password,
			m.OutputFile,
			m.Database,
			strings.Replace(m.Options, ",", " ", -1),
		)
	} else {
		return fmt.Sprintf(
			"--host %s --port %s -u %s -p%s --result-file=%s %s %s %s",
			m.Host,
			m.Port,
			m.Username,
			m.Password,
			m.OutputFile,
			m.Database,
			m.Table,
			strings.Replace(m.Options, ",", " ", -1),
		)
	}
}

// BackupMySQL backup and compress the mysql dump file
func (m *Manager) BackupMySQL(mysql *MySQL, archive string) error {
	// Get the full executable path for the editor.
	executable, err := exec.LookPath("mysqldump")

	sqlDumpFile := strings.Replace(archive, ".tar.gz", ".sql", -1)

	mysql.OutputFile = sqlDumpFile

	if err != nil {
		return err
	}

	command := strings.Split(
		fmt.Sprintf(`%s %s`, executable, mysql.DumpOptions()),
		" ",
	)

	cmd := exec.Command(command[0], command[1:]...)

	cmd.Stdin = os.Stdin
	// cmd.Stdout = os.Stdout
	// cmd.Stderr = os.Stderr

	err = cmd.Run()

	if err != nil {
		return err
	}

	return m.BackupDirectory(sqlDumpFile, archive)
}
