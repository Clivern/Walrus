// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/clivern/walrus/pkg"
)

// TestHttpGet test cases
func TestHttpGet(t *testing.T) {
	t.Run("TestHttpGet", func(t *testing.T) {
		httpClient := NewHTTPClient()
		response, error := httpClient.Get(
			context.TODO(),
			"https://httpbin.org/get",
			map[string]string{"arg1": "value1"},
			map[string]string{"X-Api-Key": "poodle-123"},
		)

		pkg.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
		pkg.Expect(t, nil, error)

		body, error := httpClient.ToString(response)

		pkg.Expect(t, true, strings.Contains(body, "value1"))
		pkg.Expect(t, true, strings.Contains(body, "arg1"))
		pkg.Expect(t, true, strings.Contains(body, "arg1=value1"))
		pkg.Expect(t, true, strings.Contains(body, "X-Api-Key"))
		pkg.Expect(t, true, strings.Contains(body, "poodle-123"))
		pkg.Expect(t, nil, error)
	})
}

// TestHttpDelete test cases
func TestHttpDelete(t *testing.T) {
	t.Run("TestHttpDelete", func(t *testing.T) {
		httpClient := NewHTTPClient()
		response, error := httpClient.Delete(
			context.TODO(),
			"https://httpbin.org/delete",
			map[string]string{"arg1": "value1"},
			map[string]string{"X-Api-Key": "poodle-123"},
		)

		pkg.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
		pkg.Expect(t, nil, error)

		body, error := httpClient.ToString(response)

		pkg.Expect(t, true, strings.Contains(body, "value1"))
		pkg.Expect(t, true, strings.Contains(body, "arg1"))
		pkg.Expect(t, true, strings.Contains(body, "arg1=value1"))
		pkg.Expect(t, true, strings.Contains(body, "X-Api-Key"))
		pkg.Expect(t, true, strings.Contains(body, "poodle-123"))
		pkg.Expect(t, nil, error)
	})
}

// TestHttpPost test cases
func TestHttpPost(t *testing.T) {
	t.Run("TestHttpPost", func(t *testing.T) {
		httpClient := NewHTTPClient()
		response, error := httpClient.Post(
			context.TODO(),
			"https://httpbin.org/post",
			`{"Username":"admin", "Password":"12345"}`,
			map[string]string{"arg1": "value1"},
			map[string]string{"X-Api-Key": "poodle-123"},
		)

		pkg.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
		pkg.Expect(t, nil, error)

		body, error := httpClient.ToString(response)

		pkg.Expect(t, true, strings.Contains(body, `"12345"`))
		pkg.Expect(t, true, strings.Contains(body, `"Username"`))
		pkg.Expect(t, true, strings.Contains(body, `"admin"`))
		pkg.Expect(t, true, strings.Contains(body, `"Password"`))
		pkg.Expect(t, true, strings.Contains(body, "value1"))
		pkg.Expect(t, true, strings.Contains(body, "arg1"))
		pkg.Expect(t, true, strings.Contains(body, "arg1=value1"))
		pkg.Expect(t, true, strings.Contains(body, "X-Api-Key"))
		pkg.Expect(t, true, strings.Contains(body, "poodle-123"))
		pkg.Expect(t, nil, error)
	})
}

// TestHttpPut test cases
func TestHttpPut(t *testing.T) {
	t.Run("TestHttpPut", func(t *testing.T) {
		httpClient := NewHTTPClient()
		response, error := httpClient.Put(
			context.TODO(),
			"https://httpbin.org/put",
			`{"Username":"admin", "Password":"12345"}`,
			map[string]string{"arg1": "value1"},
			map[string]string{"X-Api-Key": "poodle-123"},
		)

		pkg.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
		pkg.Expect(t, nil, error)

		body, error := httpClient.ToString(response)

		pkg.Expect(t, true, strings.Contains(body, `"12345"`))
		pkg.Expect(t, true, strings.Contains(body, `"Username"`))
		pkg.Expect(t, true, strings.Contains(body, `"admin"`))
		pkg.Expect(t, true, strings.Contains(body, `"Password"`))
		pkg.Expect(t, true, strings.Contains(body, "value1"))
		pkg.Expect(t, true, strings.Contains(body, "arg1"))
		pkg.Expect(t, true, strings.Contains(body, "arg1=value1"))
		pkg.Expect(t, true, strings.Contains(body, "X-Api-Key"))
		pkg.Expect(t, true, strings.Contains(body, "poodle-123"))
		pkg.Expect(t, nil, error)
	})
}

// TestHttpGetStatusCode1 test cases
func TestHttpGetStatusCode1(t *testing.T) {
	t.Run("TestHttpGetStatusCode1", func(t *testing.T) {
		httpClient := NewHTTPClient()
		response, error := httpClient.Get(
			context.TODO(),
			"https://httpbin.org/status/200",
			map[string]string{"arg1": "value1"},
			map[string]string{"X-Api-Key": "poodle-123"},
		)

		pkg.Expect(t, http.StatusOK, httpClient.GetStatusCode(response))
		pkg.Expect(t, nil, error)

		body, error := httpClient.ToString(response)

		pkg.Expect(t, "", body)
		pkg.Expect(t, nil, error)
	})
}

// TestHttpGetStatusCode2 test cases
func TestHttpGetStatusCode2(t *testing.T) {
	t.Run("TestHttpGetStatusCode2", func(t *testing.T) {
		httpClient := NewHTTPClient()
		response, error := httpClient.Get(
			context.TODO(),
			"https://httpbin.org/status/500",
			map[string]string{"arg1": "value1"},
			map[string]string{"X-Api-Key": "poodle-123"},
		)

		pkg.Expect(t, http.StatusInternalServerError, httpClient.GetStatusCode(response))
		pkg.Expect(t, nil, error)

		body, error := httpClient.ToString(response)

		pkg.Expect(t, "", body)
		pkg.Expect(t, nil, error)
	})
}

// TestHttpGetStatusCode3 test cases
func TestHttpGetStatusCode3(t *testing.T) {
	t.Run("TestHttpGetStatusCode3", func(t *testing.T) {
		httpClient := NewHTTPClient()
		response, error := httpClient.Get(
			context.TODO(),
			"https://httpbin.org/status/404",
			map[string]string{"arg1": "value1"},
			map[string]string{"X-Api-Key": "poodle-123"},
		)

		pkg.Expect(t, http.StatusNotFound, httpClient.GetStatusCode(response))
		pkg.Expect(t, nil, error)

		body, error := httpClient.ToString(response)

		pkg.Expect(t, "", body)
		pkg.Expect(t, nil, error)
	})
}

// TestHttpGetStatusCode4 test cases
func TestHttpGetStatusCode4(t *testing.T) {
	t.Run("TestHttpGetStatusCode4", func(t *testing.T) {
		httpClient := NewHTTPClient()
		response, error := httpClient.Get(
			context.TODO(),
			"https://httpbin.org/status/201",
			map[string]string{"arg1": "value1"},
			map[string]string{"X-Api-Key": "poodle-123"},
		)

		pkg.Expect(t, http.StatusCreated, httpClient.GetStatusCode(response))
		pkg.Expect(t, nil, error)

		body, error := httpClient.ToString(response)

		pkg.Expect(t, "", body)
		pkg.Expect(t, nil, error)
	})
}

// TestBuildParameters test cases
func TestBuildParameters(t *testing.T) {
	t.Run("TestBuildParameters", func(t *testing.T) {
		httpClient := NewHTTPClient()
		url, error := httpClient.BuildParameters("http://127.0.0.1", map[string]string{"arg1": "value1"})

		pkg.Expect(t, "http://127.0.0.1?arg1=value1", url)
		pkg.Expect(t, nil, error)
	})
}

// TestBuildData test cases
func TestBuildData(t *testing.T) {
	t.Run("TestBuildData", func(t *testing.T) {
		httpClient := NewHTTPClient()
		pkg.Expect(t, httpClient.BuildData(map[string]string{}), "")
		pkg.Expect(t, httpClient.BuildData(map[string]string{"arg1": "value1"}), "arg1=value1")
	})
}
