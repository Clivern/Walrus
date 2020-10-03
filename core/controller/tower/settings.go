// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"io/ioutil"
	"net/http"

	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetSettings controller
func GetSettings(c *gin.Context) {
	db := service.Database{}

	err := db.AutoConnect()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Error(`Failure while connecting database`)

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	c.JSON(http.StatusOK, gin.H{
		"options": db.GetOptions(),
	})
}

// UpdateSettings controller
func UpdateSettings(c *gin.Context) {
	var options model.Options

	db := service.Database{}

	x, _ := ioutil.ReadAll(c.Request.Body)

	options.LoadFromJSON(x)

	err := db.AutoConnect()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Error(`Failure while connecting database`)

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	c.Status(http.StatusOK)
	return
}
