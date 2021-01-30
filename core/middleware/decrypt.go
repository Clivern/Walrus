// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package middleware

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/clivern/walrus/core/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Decrypt middleware
func Decrypt() gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte
		var decBodyBytes []byte
		var err error

		encKey := ""
		path := c.Request.URL.Path
		encRoute := false

		if strings.Contains(path, "/bootstrap_agent") ||
			strings.Contains(path, "/agent_heartbeat") ||
			strings.Contains(path, "/agent_postback") ||
			strings.Contains(path, "/process") {
			encRoute = true
		}

		// All unencrypted requests if env is dev
		if viper.GetString(fmt.Sprintf("%s.mode", viper.GetString("role"))) == "dev" {
			encRoute = false
		}

		if viper.GetString("role") == "tower" {
			encKey = viper.GetString("tower.api.encryptionKey")
		} else if viper.GetString("role") == "agent" {
			encKey = viper.GetString("agent.tower.encryptionKey")
		}

		encrypted := c.GetHeader("x-encrypted-request")

		if encKey == "" {
			log.Error(`Encryption Key is missing`)
			return
		}

		if encrypted == "" && !encRoute {
			return
		}

		// Workaround for issue https://github.com/gin-gonic/gin/issues/1651
		if c.Request.Body != nil {
			bodyBytes, _ = ioutil.ReadAll(c.Request.Body)
			decBodyBytes, err = util.Decrypt(bodyBytes, encKey)

			if err != nil {
				log.WithFields(log.Fields{
					"correlation_id": c.GetHeader("x-correlation-id"),
					"http_method":    c.Request.Method,
					"http_path":      c.Request.URL.Path,
					"request_body":   string(bodyBytes),
					"error":          err.Error(),
				}).Info(`Invalid encrypted request`)

				c.AbortWithStatus(http.StatusBadRequest)
			}
		}

		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(decBodyBytes))

		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"http_method":    c.Request.Method,
			"http_path":      c.Request.URL.Path,
		}).Info("Request body decrypted")
	}
}
