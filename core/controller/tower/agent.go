// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/module"
	"github.com/clivern/walrus/core/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// AgentBootstrap controller
func AgentBootstrap(c *gin.Context) {
	var bootstrapRequest module.BootstrapRequest

	data, _ := ioutil.ReadAll(c.Request.Body)

	err := util.LoadFromJSON(&bootstrapRequest, data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Error! Invalid request",
		})
		return
	}

	db := driver.NewEtcdDriver()

	err = db.Connect()

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

	agentStore := model.NewAgentStore(db)
	hostStore := model.NewHostStore(db)

	_, err = hostStore.GetHost(bootstrapRequest.Hostname)

	// Create host if missing
	if err != nil {
		err = hostStore.CreateHost(model.HostData{
			Hostname: bootstrapRequest.Hostname,
		})

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
	}

	err = agentStore.CreateAgent(model.AgentData{
		ID:              bootstrapRequest.AgentID,
		URL:             bootstrapRequest.AgentURL,
		Hostname:        bootstrapRequest.Hostname,
		APIKey:          bootstrapRequest.AgentAPIKey,
		CreatedAt:       time.Now().Unix(),
		UpdatedAt:       time.Now().Unix(),
		LastStatusCheck: time.Now().Unix(),
	})

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

	c.Status(http.StatusOK)
	return
}

// AgentHeartbeat controller
func AgentHeartbeat(c *gin.Context) {
	var heartbeatRequest module.HeartbeatRequest

	data, _ := ioutil.ReadAll(c.Request.Body)

	err := util.LoadFromJSON(&heartbeatRequest, data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Error! Invalid request",
		})
		return
	}

	db := driver.NewEtcdDriver()

	err = db.Connect()

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

	agent := model.NewAgentStore(db)

	agentData, err := agent.GetAgent(heartbeatRequest.Hostname, heartbeatRequest.AgentID)

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

	agentData.Status = heartbeatRequest.Status
	agentData.LastStatusCheck = time.Now().Unix()

	err = agent.UpdateAgent(*agentData)

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

	c.Status(http.StatusOK)
	return
}

// AgentPostback controller
func AgentPostback(c *gin.Context) {
	var postbackRequest module.PostbackRequest

	data, _ := ioutil.ReadAll(c.Request.Body)

	err := util.LoadFromJSON(&postbackRequest, data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Error! Invalid request",
		})
		return
	}

	if postbackRequest.Status == model.PendingStatus {
		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Invalid job status received",
		})
		return
	}

	db := driver.NewEtcdDriver()

	err = db.Connect()

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

	jobStore := model.NewJobStore(db)
	cronStore := model.NewCronStore(db)

	// Delete the Job
	_, err = jobStore.DeleteRecord(postbackRequest.Hostname, postbackRequest.JobID)

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

	// Increment the cron jobs status
	cronData, err := cronStore.GetRecord(postbackRequest.Hostname, postbackRequest.CronID)

	// If cron missing, skip
	if err != nil && strings.Contains(err.Error(), "Unable to find") {
		c.Status(http.StatusOK)
		return
	}

	if postbackRequest.Status == model.FailedStatus {
		cronData.FailedJobs++
	} else if postbackRequest.Status == model.SuccessStatus {
		cronData.SuccessJobs++
	}

	err = cronStore.UpdateRecord(*cronData)

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

	c.Status(http.StatusOK)
	return
}
