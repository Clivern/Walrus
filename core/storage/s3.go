// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package storage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3 type
type S3 struct {
	config *aws.Config
}

// NewS3Client creates a new s3 instance
func NewS3Client(key, secret, endpoint, region string) *S3 {
	return &S3{
		config: &aws.Config{
			Credentials: credentials.NewStaticCredentials(key, secret, ""),
			Endpoint:    aws.String(endpoint),
			Region:      aws.String(region),
		},
	}
}

// CreateBucket creates a bucket
func (s *S3) CreateBucket(bucket string) error {
	newSession := session.New(s.config)
	s3Client := s3.New(newSession)

	params := &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	}

	_, err := s3Client.CreateBucket(params)

	if err != nil {
		return err
	}

	return nil
}

// UploadFile upload a file to s3 bucket
func (s *S3) UploadFile(bucket, localPath, remotePath string, includeChecksum bool) error {

	file, err := os.Open(localPath)

	if err != nil {
		return err
	}

	defer file.Close()

	newSession := session.New(s.config)

	uploader := s3manager.NewUploader(newSession)

	// Upload the file's body to S3 bucket
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(remotePath),
		Body:   file,
	})

	if err != nil {
		return err
	}

	if !includeChecksum {
		return nil
	}

	h := sha256.New()

	if _, err := io.Copy(h, file); err != nil {
		return err
	}

	checksumContent := fmt.Sprintf("SHA256 CheckSum: %x\n", h.Sum(nil))

	checksumPath := fmt.Sprintf("%s-checksum.txt",
		strings.TrimSuffix(remotePath, filepath.Ext(remotePath)),
	)

	tmpFile, err := ioutil.TempFile(os.TempDir(), "walrus-")

	if err != nil {
		return err
	}

	// Remember to clean up the file afterwards
	defer os.Remove(tmpFile.Name())

	if _, err = tmpFile.Write([]byte(checksumContent)); err != nil {
		return err
	}

	if err := tmpFile.Close(); err != nil {
		return err
	}

	sumFile, err := os.Open(tmpFile.Name())

	if err != nil {
		return err
	}

	defer sumFile.Close()

	// Upload the file's body to S3 bucket
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(checksumPath),
		Body:   sumFile,
	})

	if err != nil {
		return err
	}

	return nil
}
