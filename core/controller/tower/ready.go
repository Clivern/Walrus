// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"net/http"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/module"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Ready controller
func Ready(c *gin.Context) {
	db := driver.NewEtcdDriver()

	err := db.Connect()

	if err != nil || !db.IsConnected() {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"status":         "NotOk",
		}).Info(`Ready check`)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "NotOk",
		})

		return
	}

	defer db.Close()

	// Check if there is a leader tower
	tower := module.NewTower(db)
	hasLeader, err := tower.HasLeader()

	if !hasLeader || err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"status":         "NotOk",
		}).Info(`Ready check`)

		c.JSON(http.StatusInternalServerError, gin.H{
			"status": "NotOk",
		})
		return
	}

	log.WithFields(log.Fields{
		"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		"status":         "ok",
	}).Info(`Ready check`)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
