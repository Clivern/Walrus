// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package middleware

import (
	"net/http"
	"strings"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Auth middleware
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.Request.URL.Path
		method := c.Request.Method
		reqAPIKey := c.GetHeader("x-api-key")
		clientID := c.GetHeader("x-client-id")
		clientEmail := c.GetHeader("x-user-email")

		apiKey := ""

		if viper.GetString("role") == "tower" {
			apiKey = viper.GetString("tower.api.key")
		} else if viper.GetString("role") == "agent" {
			apiKey = viper.GetString("agent.api.key")
		}

		// Skip if endpoint not protected API
		if !strings.Contains(path, "/api/") {
			return
		}

		// Validate API Call
		if clientID != "dashboard" {
			if apiKey != "" && apiKey != reqAPIKey {
				log.WithFields(log.Fields{
					"correlation_id":  c.GetHeader("x-correlation-id"),
					"http_method":     method,
					"http_path":       path,
					"request_api_key": reqAPIKey,
				}).Info(`Unauthorized access`)

				c.AbortWithStatus(http.StatusUnauthorized)
			}

			return
		}

		// Validate Logged User Request
		if !util.IsEmailValid(clientEmail) || util.IsEmpty(reqAPIKey) {
			log.WithFields(log.Fields{
				"correlation_id": c.GetHeader("x-correlation-id"),
				"http_method":    method,
				"http_path":      path,
				"api_key":        reqAPIKey,
				"client_email":   clientEmail,
			}).Info(`Unauthorized access`)

			c.AbortWithStatus(http.StatusUnauthorized)
		}

		db := driver.NewEtcdDriver()

		err := db.Connect()

		if err != nil {
			log.WithFields(log.Fields{
				"correlation_id": c.GetHeader("x-correlation-id"),
				"error":          err.Error(),
			}).Error("Internal server error")

			c.JSON(http.StatusInternalServerError, gin.H{
				"correlationID": c.GetHeader("x-correlation-id"),
				"errorMessage":  "Internal server error",
			})
			return
		}

		defer db.Close()

		userStore := model.NewUserStore(db)

		user, err := userStore.GetUserByEmail(clientEmail)

		if err != nil || user.APIKey != reqAPIKey {
			log.WithFields(log.Fields{
				"correlation_id": c.GetHeader("x-correlation-id"),
				"http_method":    method,
				"http_path":      path,
				"api_key":        reqAPIKey,
				"client_email":   clientEmail,
			}).Info(`Unauthorized access`)

			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
