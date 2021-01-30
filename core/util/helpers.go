// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// InArray check if value is on array
func InArray(val interface{}, array interface{}) bool {
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				return true
			}
		}
	}

	return false
}

// GenerateUUID4 create a UUID
func GenerateUUID4() string {
	u := uuid.Must(uuid.NewV4(), nil)
	return u.String()
}

// ReadFile get the file content
func ReadFile(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// FilterFiles filters files list based on specific sub-strings
func FilterFiles(files, filters []string) []string {
	var filteredFiles []string

	for _, file := range files {
		ok := true
		for _, filter := range filters {

			ok = ok && strings.Contains(file, filter)
		}
		if ok {
			filteredFiles = append(filteredFiles, file)
		}
	}

	return filteredFiles
}

// Unset remove element at position i
func Unset(a []string, i int) []string {
	a[i] = a[len(a)-1]
	a[len(a)-1] = ""
	return a[:len(a)-1]
}

// EnsureTrailingSlash ensure there is a trailing slash
func EnsureTrailingSlash(dir string) string {
	return fmt.Sprintf(
		"%s%s",
		strings.TrimRight(dir, string(os.PathSeparator)),
		string(os.PathSeparator),
	)
}

// RemoveTrailingSlash removes any trailing slash
func RemoveTrailingSlash(dir string) string {
	return strings.TrimRight(dir, string(os.PathSeparator))
}

// RemoveStartingSlash removes any starting slash
func RemoveStartingSlash(dir string) string {
	return strings.TrimLeft(dir, string(os.PathSeparator))
}

// ClearDir removes all files and sub dirs
func ClearDir(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}
	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			return err
		}
	}
	return nil
}

// StoreFile stores a file content
func StoreFile(path, content string) error {
	dir := filepath.Dir(path)

	err := os.MkdirAll(dir, 0775)

	if err != nil {
		return err
	}

	f, err := os.Create(path)

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(content)

	return err
}

// PathExists reports whether the path exists
func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// FileExists reports whether the named file exists
func FileExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsRegular() {
			return true
		}
	}
	return false
}

// DirExists reports whether the dir exists
func DirExists(path string) bool {
	if fi, err := os.Stat(path); err == nil {
		if fi.Mode().IsDir() {
			return true
		}
	}
	return false
}

// EnsureDir ensures that directory exists
func EnsureDir(dirName string, mode int) (bool, error) {
	err := os.MkdirAll(dirName, os.FileMode(mode))

	if err == nil || os.IsExist(err) {
		return true, nil
	}
	return false, err
}

// DeleteFile deletes a file
func DeleteFile(path string) error {
	return os.Remove(path)
}

// Encrypt encrypt data
func Encrypt(data []byte, passphrase string) ([]byte, error) {
	var result []byte

	block, _ := aes.NewCipher([]byte(createHash(passphrase)))
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return result, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return result, err
	}

	result = gcm.Seal(nonce, nonce, data, nil)

	return result, nil
}

// Decrypt decrypts data
func Decrypt(data []byte, passphrase string) ([]byte, error) {
	var result []byte

	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)

	if err != nil {
		panic(err.Error())
	}

	gcm, err := cipher.NewGCM(block)

	if err != nil {
		return result, err
	}

	nonceSize := gcm.NonceSize()

	if nonceSize > len(data)-1 {
		return result, fmt.Errorf("Invalid encrypted text")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	result, err = gcm.Open(nil, nonce, ciphertext, nil)

	if err != nil {
		return result, err
	}

	return result, nil
}

// createHash creates a hash
func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

// LoadFromJSON update object from json
func LoadFromJSON(item interface{}, data []byte) error {
	err := json.Unmarshal(data, &item)
	if err != nil {
		return err
	}

	return nil
}

// ConvertToJSON convert object to json
func ConvertToJSON(item interface{}) (string, error) {
	data, err := json.Marshal(&item)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// HashPassword hashes the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword(
		[]byte(password),
		14,
	)

	return string(bytes), err
}

// CheckPasswordHash validate password with a hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hash),
		[]byte(password),
	)

	return err == nil
}

// GetHostname gets the hostname
func GetHostname() (string, error) {
	hostname, err := os.Hostname()

	if err != nil {
		return "", err
	}

	return strings.ToLower(hostname), nil
}

// ValidRelPath checks for path traversal and correct forward slashes
func ValidRelPath(path string) bool {
	if path == "" || strings.Contains(path, `\`) || strings.HasPrefix(path, "/") || strings.Contains(path, "../") {
		return false
	}

	return true
}

// DeleteDir deletes a dir
func DeleteDir(dir string) bool {
	if err := os.RemoveAll(dir); err == nil {
		return true
	}

	return false
}

// IsEmailValid checks if the email provided passes the required structure
// and length test. It also checks the domain has a valid MX record.
func IsEmailValid(email string) bool {
	if len(email) < 3 && len(email) > 254 {
		return false
	}

	var emailRegex = regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)

	if !emailRegex.MatchString(email) {
		return false
	}

	parts := strings.Split(email, "@")

	mx, err := net.LookupMX(parts[1])

	if err != nil || len(mx) == 0 {
		return false
	}

	return true
}

// IsEmpty validate if string is empty or not
func IsEmpty(item string) bool {
	if strings.TrimSpace(item) == "" {
		return true
	}
	return false
}
