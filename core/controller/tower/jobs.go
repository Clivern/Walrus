// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"fmt"
	"net/http"

	"github.com/clivern/walrus/core/model"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// GetJobs controller
func GetJobs(c *gin.Context) {
	db := model.Database{}

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
		"jobs": db.GetJobs(),
	})
}

// GetJob controller
func GetJob(c *gin.Context) {
	uuid := c.Param("uuid")

	db := model.Database{}

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

	job := db.GetJobByUUID(uuid)

	if job.ID < 1 {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"job_uuid":       uuid,
		}).Info(fmt.Sprintf(`Job not found`))

		c.Status(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		"job_uuid":       uuid,
	}).Info(`Retrieve a job`)

	c.JSON(http.StatusOK, gin.H{
		"id":        job.ID,
		"uuid":      job.UUID,
		"status":    job.Status,
		"type":      job.Type,
		"runAt":     job.RunAt,
		"createdAt": job.CreatedAt,
		"updatedAt": job.UpdatedAt,
	})
}

// DeleteJob controller
func DeleteJob(c *gin.Context) {
	uuid := c.Param("uuid")

	db := model.Database{}

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

	job := db.GetJobByUUID(uuid)

	if job.ID < 1 {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"job_uuid":       uuid,
		}).Info(`Job not found`)

		c.Status(http.StatusNotFound)
		return
	}

	log.WithFields(log.Fields{
		"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
		"job_uuid":       uuid,
	}).Info(`Deleting a job`)

	db.DeleteJobByID(job.ID)

	c.Status(http.StatusNoContent)
	return
}
