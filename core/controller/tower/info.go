// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"net/http"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Info controller
func Info(c *gin.Context) {
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

	optionStore := model.NewOptionStore(db)

	_, err = optionStore.GetOptionByKey("is_installed")

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"setupStatus": false,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"setupStatus": true,
	})
}
