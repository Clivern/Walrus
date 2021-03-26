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
	"time"

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
			Credentials:      credentials.NewStaticCredentials(key, secret, ""),
			Endpoint:         aws.String(endpoint),
			Region:           aws.String(region),
			S3ForcePathStyle: aws.Bool(true),
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
		strings.TrimSuffix(remotePath, ".tar.gz"),
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

// ListFiles lists files
func (s *S3) ListFiles(bucket, prefix string) ([]string, error) {
	newSession := session.New(s.config)
	s3Client := s3.New(newSession)

	result := []string{}

	err := s3Client.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}, func(p *s3.ListObjectsOutput, last bool) (shouldContinue bool) {
		for _, obj := range p.Contents {
			result = append(result, fmt.Sprintf("%s", *obj.Key))
		}

		return true
	})

	if err != nil {
		return result, err
	}

	return result, nil
}

// DeleteFile delete a file
func (s *S3) DeleteFile(bucket, file string) error {
	newSession := session.New(s.config)
	s3Client := s3.New(newSession)

	_, err := s3Client.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(file),
	})

	if err != nil {
		return err
	}

	return nil
}

// CleanupOld apply a retention policy over a certain path
func (s *S3) CleanupOld(bucket, path string, beforeDays int) (int, error) {
	files, err := s.ListFiles(bucket, path)

	count := 0

	deleteBefore := time.Now().Add(time.Duration(-1*beforeDays*24) * time.Hour)

	if err != nil {
		return count, err
	}

	for _, file := range files {
		fileName := filepath.Base(file)
		fileName = strings.TrimSuffix(fileName, ".tar.gz")
		fileName = strings.TrimSuffix(fileName, "-checksum.txt")
		fileTime, _ := time.Parse("2006-01-02_15-04-05", fileName)

		if fileTime.Before(deleteBefore) {
			count++
			s.DeleteFile(bucket, file)
		}
	}

	return count, nil
}
