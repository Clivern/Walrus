// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"fmt"
	"net/http"

	"github.com/clivern/walrus/core/service"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetHosts controller
func GetHosts(c *gin.Context) {
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
		"services": db.GetHosts(),
	})
}

// GetHost controller
func GetHost(c *gin.Context) {
	uuid := c.Param("hostId")

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

	serv := db.GetHostByUUID(uuid)

	if serv.ID < 1 {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"service_uuid":   uuid,
		}).Info(fmt.Sprintf(`Host not found`))

		c.Status(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		"service_uuid":   uuid,
	}).Info(`Retrieve a service`)

	c.JSON(http.StatusOK, gin.H{
		"id":        serv.ID,
		"uuid":      serv.UUID,
		"status":    serv.Status,
		"type":      serv.Type,
		"destroyAt": serv.DestroyAt,
		"configs":   serv.Configs,
		"createdAt": serv.CreatedAt,
		"updatedAt": serv.UpdatedAt,
	})
}

// DeleteHost controller
func DeleteHost(c *gin.Context) {
	uuid := c.Param("hostId")

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

	serv := db.GetHostByUUID(uuid)

	if serv.ID < 1 {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"service_uuid":   uuid,
		}).Info(`Host not found`)

		c.Status(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		"service_uuid":   uuid,
	}).Info(`Deleting a service`)

	db.DeleteHostByID(serv.ID)

	c.Status(http.StatusNoContent)
	return
}

// UpdateHost controller
func UpdateHost(c *gin.Context) {
	c.Status(http.StatusAccepted)
	return
}
