// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package agent

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/clivern/walrus/core/module"
	"github.com/clivern/walrus/core/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Workers type
type Workers struct {
	broadcast chan module.BackupMessage
}

// NewWorkers get a new workers instance
func NewWorkers() *Workers {
	result := new(Workers)
	result.broadcast = make(chan module.BackupMessage, viper.GetInt(
		fmt.Sprintf("%s.workers.buffer", viper.GetString("role")),
	))

	return result
}

// BroadcastRequest sends a tower request to workers
func (w *Workers) BroadcastRequest(c *gin.Context, rawBody []byte) {
	message := &module.BackupMessage{}

	err := util.LoadFromJSON(message, rawBody)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error(`Invalid backup message`)

		c.JSON(http.StatusInternalServerError, gin.H{
			"errorMessage": "Internal server error",
		})

		return
	}

	log.WithFields(log.Fields{
		"correlation_id": message.CorrelationID,
		"message":        message,
	}).Info(`Incoming request`)

	w.broadcast <- *message

	c.Status(http.StatusAccepted)
	return
}

// HandleWorkload handles all incoming requests from tower
func (w *Workers) HandleWorkload() <-chan module.BackupMessage {
	notifyChannel := make(chan module.BackupMessage)

	go func() {
		wg := &sync.WaitGroup{}

		for t := 0; t < viper.GetInt(fmt.Sprintf("%s.workers.count", viper.GetString("role"))); t++ {
			wg.Add(1)
			go w.ProcessAction(notifyChannel, wg)
		}

		wg.Wait()

		close(notifyChannel)
	}()

	return notifyChannel
}

// ProcessAction process incoming request from the tower
func (w *Workers) ProcessAction(notifyChannel chan<- module.BackupMessage, wg *sync.WaitGroup) {
	for message := range w.broadcast {
		log.WithFields(log.Fields{
			"correlation_id": message.CorrelationID,
			"message":        message,
		}).Info(`Worker received a new message`)

		// ~

		log.WithFields(log.Fields{
			"correlation_id": message.CorrelationID,
			"message":        message,
		}).Info(`Worker finished processing`)

		notifyChannel <- message
	}

	wg.Done()
}

// NotifyTower notifies tower
func (w *Workers) NotifyTower(notifyChannel <-chan module.BackupMessage) {
	for message := range notifyChannel {
		log.WithFields(log.Fields{
			"correlation_id": message.CorrelationID,
			"message":        message,
		}).Info(`Worker finished processing`)
	}
}
