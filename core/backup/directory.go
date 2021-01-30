// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package backup

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/clivern/walrus/core/util"
)

// Directory struct
type Directory struct {
}

// Backup creates a backup of any directory
// directory must be absolute path like /etc/app and archive same /etc/app.tar.gzip"
func (d *Directory) Backup(directory, archive string) error {
	var buf bytes.Buffer

	if !util.DirExists(directory) {
		return fmt.Errorf("Unable to find directory %s", directory)
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
		return err
	}

	// Produce tar
	if err := tarWriter.Close(); err != nil {
		return err
	}
	// Produce gzip
	if err := gzipWriter.Close(); err != nil {
		return err
	}

	targetTar, err := os.OpenFile(archive, os.O_CREATE|os.O_RDWR, os.FileMode(0755))

	if err != nil {
		return err
	}

	if _, err := io.Copy(targetTar, &buf); err != nil {
		return err
	}

	return nil
}

// Restore uncompress .tar.gzip archive
func (d *Directory) Restore(archive, destination string) error {

	if !util.FileExists(archive) {
		return fmt.Errorf("Unable to find archive %s", archive)
	}

	file, err := os.Open(archive)

	if err != nil {
		return err
	}

	defer file.Close()

	zr, err := gzip.NewReader(file)

	if err != nil {
		return err
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
			return err
		}

		target := strings.Replace(
			header.Name,
			util.EnsureTrailingSlash(destination),
			"",
			-1,
		)

		target = util.RemoveTrailingSlash(target)
		target = util.RemoveStartingSlash(target)

		// Validate name against path traversal
		if !util.ValidRelPath(target) {
			return fmt.Errorf("tar contained invalid name error %s", target)
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
					return err
				}
			}
		// If it's a file create it (with same permission)
		case tar.TypeReg:
			fileToWrite, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// Copy over contents
			if _, err := io.Copy(fileToWrite, tr); err != nil {
				return err
			}

			// Manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			fileToWrite.Close()
		}
	}

	return nil
}
