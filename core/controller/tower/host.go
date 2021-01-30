// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Host type
type Host struct {
	ID           string    `json:"id"`
	Hostname     string    `json:"hostname"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	OnlineAgents int       `json:"onlineAgents"`
}

// HostCron type
type HostCron struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Request      model.Request `json:"request"`
	Interval     int           `json:"interval"`
	IntervalType string        `json:"intervalType"`
	LastRun      time.Time     `json:"lastRun"`
	SuccessJobs  int           `json:"successJobs"`
	FailedJobs   int           `json:"failedJobs"`
	PendingJobs  int           `json:"pendingJobs"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

// GetHosts controller
func GetHosts(c *gin.Context) {
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

	hostStore := model.NewHostStore(db)
	agentStore := model.NewAgentStore(db)
	data, err := hostStore.GetHosts()

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

	var hosts []Host

	for _, v := range data {
		count, _ := agentStore.CountOnlineAgents(v.Hostname)

		hosts = append(hosts, Host{
			ID:           v.ID,
			Hostname:     v.Hostname,
			CreatedAt:    time.Unix(v.CreatedAt, 0),
			UpdatedAt:    time.Unix(v.UpdatedAt, 0),
			OnlineAgents: count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"hosts": hosts,
	})
}

// GetHost controller
func GetHost(c *gin.Context) {

	hostname := c.Param("hostname")
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

	hostStore := model.NewHostStore(db)
	agentStore := model.NewAgentStore(db)
	hostData, err := hostStore.GetHost(hostname)

	if err != nil && strings.Contains(err.Error(), "Unable to find") {
		c.JSON(http.StatusNotFound, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  fmt.Sprintf("Unable to fine host: %s", hostname),
		})
		return
	}

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusNotFound, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	count, _ := agentStore.CountOnlineAgents(hostData.Hostname)

	c.JSON(http.StatusOK, gin.H{
		"id":           hostData.ID,
		"hostname":     hostData.Hostname,
		"createdAt":    time.Unix(hostData.CreatedAt, 0),
		"onlineAgents": count,
		"updatedAt":    time.Unix(hostData.UpdatedAt, 0),
	})
}

// GetHostCrons controller
func GetHostCrons(c *gin.Context) {
	db := driver.NewEtcdDriver()

	hostname := c.Param("hostname")
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

	cronStore := model.NewCronStore(db)
	jobStore := model.NewJobStore(db)

	data, err := cronStore.GetHostCrons(hostname)

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

	var crons []HostCron

	for _, v := range data {
		pendingCount, _ := jobStore.CountHostJobs(hostname, model.PendingStatus)

		crons = append(crons, HostCron{
			ID:           v.ID,
			Name:         v.Name,
			Request:      v.Request,
			Interval:     v.Interval,
			IntervalType: v.IntervalType,
			SuccessJobs:  v.SuccessJobs,
			FailedJobs:   v.FailedJobs,
			PendingJobs:  pendingCount,
			LastRun:      time.Unix(v.LastRun, 0),
			CreatedAt:    time.Unix(v.CreatedAt, 0),
			UpdatedAt:    time.Unix(v.UpdatedAt, 0),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"crons": crons,
	})
}

// CreateHostCron controller
func CreateHostCron(c *gin.Context) {

	var inputs map[string]string

	data, _ := ioutil.ReadAll(c.Request.Body)

	err := util.LoadFromJSON(&inputs, data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Error! Invalid request",
		})
		return
	}

	db := driver.NewEtcdDriver()

	hostname := c.Param("hostname")
	cronID := util.GenerateUUID4()

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

	cronStore := model.NewCronStore(db)

	interval, _ := strconv.Atoi(inputs["interval"])
	retention, _ := strconv.Atoi(inputs["retention"])

	err = cronStore.CreateRecord(model.CronRecord{
		ID:       cronID,
		Hostname: hostname,
		Name:     inputs["name"],
		Request: model.Request{
			Type:          model.BackupDirectory,
			Directory:     inputs["directory"],
			RetentionDays: retention,
		},
		Interval:     interval,
		IntervalType: inputs["intervalType"],
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

	c.JSON(http.StatusCreated, gin.H{
		"id": cronID,
	})
}

// UpdateHostCron controller
func UpdateHostCron(c *gin.Context) {

	var inputs map[string]string

	data, _ := ioutil.ReadAll(c.Request.Body)

	err := util.LoadFromJSON(&inputs, data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Error! Invalid request",
		})
		return
	}

	db := driver.NewEtcdDriver()

	hostname := c.Param("hostname")
	cronID := c.Param("cronId")

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

	cronStore := model.NewCronStore(db)

	cronData, err := cronStore.GetRecord(hostname, cronID)

	if err != nil && strings.Contains(err.Error(), "Unable to find") {
		c.JSON(http.StatusNotFound, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  fmt.Sprintf("Unable to fine host: %s", hostname),
		})
		return
	}

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

	interval, _ := strconv.Atoi(inputs["interval"])
	retention, _ := strconv.Atoi(inputs["retention"])

	cronData.Name = inputs["name"]
	cronData.Request.Directory = inputs["directory"]
	cronData.Request.RetentionDays = retention
	cronData.Interval = interval
	cronData.IntervalType = inputs["intervalType"]

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

	c.JSON(http.StatusCreated, gin.H{
		"id": cronID,
	})
}

// DeleteHostCron controller
func DeleteHostCron(c *gin.Context) {
	db := driver.NewEtcdDriver()

	hostname := c.Param("hostname")
	cronID := c.Param("cronId")
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

	cronStore := model.NewCronStore(db)
	_, err = cronStore.DeleteRecord(hostname, cronID)

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

	c.Status(http.StatusNoContent)
	return
}

// GetHostCron controller
func GetHostCron(c *gin.Context) {
	db := driver.NewEtcdDriver()

	hostname := c.Param("hostname")
	cronID := c.Param("cronId")
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

	cronStore := model.NewCronStore(db)
	cronData, err := cronStore.GetRecord(hostname, cronID)

	if err != nil && strings.Contains(err.Error(), "Unable to find") {
		c.JSON(http.StatusNotFound, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  fmt.Sprintf("Unable to fine host: %s", hostname),
		})
		return
	}

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

	c.JSON(http.StatusOK, gin.H{
		"id":           cronData.ID,
		"name":         cronData.Name,
		"hostname":     cronData.Hostname,
		"interval":     cronData.Interval,
		"intervalType": cronData.IntervalType,
		"request": gin.H{
			"type":          cronData.Request.Type,
			"directory":     cronData.Request.Directory,
			"retentionDays": cronData.Request.RetentionDays,
		},
		"createdAt": cronData.CreatedAt,
		"lastRun":   cronData.LastRun,
		"updatedAt": cronData.UpdatedAt,
	})
}

// DeleteHost controller
func DeleteHost(c *gin.Context) {
	db := driver.NewEtcdDriver()

	hostname := c.Param("hostname")
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

	hostStore := model.NewHostStore(db)
	_, err = hostStore.DeleteHost(hostname)

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

	c.Status(http.StatusNoContent)
	return
}
