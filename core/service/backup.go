// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Backup struct
type Backup struct {
}

// BackupDirectory creates a backup of any directory
// directory must be absolute path like /etc/app and archive same /etc/app.tar.gzip"
func (b *Backup) BackupDirectory(directory, archive string) (bool, error) {
	var buf bytes.Buffer

	if !DirExists(directory) {
		return false, fmt.Errorf("Unable to find directory %s", directory)
	}

	gzipWriter := gzip.NewWriter(&buf)
	tarWriter := tar.NewWriter(gzipWriter)

	// Walk through every file in the folder
	err := filepath.Walk(directory, func(file string, fi os.FileInfo, err error) error {
		// Generate tar header
		header, err := tar.FileInfoHeader(fi, file)

		if err != nil {
			return err
		}

		header.Name = filepath.ToSlash(file)

		// Write header
		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		// If not a dir, write file content
		if !fi.IsDir() {
			data, err := os.Open(file)

			if err != nil {
				return err
			}

			if _, err := io.Copy(tarWriter, data); err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	// Produce tar
	if err := tarWriter.Close(); err != nil {
		return false, err
	}
	// Produce gzip
	if err := gzipWriter.Close(); err != nil {
		return false, err
	}

	targetTar, err := os.OpenFile(archive, os.O_CREATE|os.O_RDWR, os.FileMode(0755))

	if err != nil {
		return false, err
	}

	if _, err := io.Copy(targetTar, &buf); err != nil {
		return false, err
	}

	return true, nil
}

// RestoreDirectory uncompress .tar.gzip archive
func (b *Backup) RestoreDirectory(archive, destination string) (bool, error) {

	if !FileExists(archive) {
		return false, fmt.Errorf("Unable to find archive %s", archive)
	}

	file, err := os.Open(archive)

	if err != nil {
		return false, err
	}

	defer file.Close()

	zr, err := gzip.NewReader(file)

	if err != nil {
		return false, err
	}

	// Untar
	tr := tar.NewReader(zr)

	// Uncompress each element
	for {
		header, err := tr.Next()

		if err == io.EOF {
			break // End of archive
		}

		if err != nil {
			return false, err
		}

		target := strings.Replace(
			header.Name,
			EnsureTrailingSlash(destination),
			"",
			-1,
		)

		target = RemoveTrailingSlash(target)
		target = RemoveStartingSlash(target)

		// Validate name against path traversal
		if !ValidRelPath(target) {
			return false, fmt.Errorf("tar contained invalid name error %q\n", target)
		}

		// Add destination + re-format slashes according to system
		target = filepath.Join(destination, target)
		// if no join is needed, replace with ToSlash:
		// target = filepath.ToSlash(header.Name)

		// Check the type
		switch header.Typeflag {

		// If its a dir and it doesn't exist create it (with 0755 permission)
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return false, err
				}
			}
		// If it's a file create it (with same permission)
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return false, err
			}

			// Copy over contents
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return false, err
			}

			// Manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			fileToWrite.Close()
		}
	}

	return true, nil
}

// BackupMySQLDatabase ...
func (b *Backup) BackupMySQLDatabase() (bool, error) {
	return true, nil
}

// RestoreMySQLDatabase ...
func (b *Backup) RestoreMySQLDatabase() (bool, error) {
	return true, nil
}

// BackupPostgresqlDatabase ...
func (b *Backup) BackupPostgresqlDatabase() (bool, error) {
	return true, nil
}

// RestorePostgresqlDatabase ...
func (b *Backup) RestorePostgresqlDatabase() (bool, error) {
	return true, nil
}
