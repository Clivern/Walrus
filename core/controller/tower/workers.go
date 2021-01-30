// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"fmt"
	"sync"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/module"
	"github.com/clivern/walrus/core/service"
	"github.com/clivern/walrus/core/util"

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

// BroadcastRequests sends a tower request to workers
func (w *Workers) BroadcastRequests() {

	db := driver.NewEtcdDriver()
	err := db.Connect()

	if err != nil {
		panic(fmt.Sprintf(
			"Error while connecting to etcd: %s",
			err.Error(),
		))
	}

	cronStore := model.NewCronStore(db)
	jobStore := model.NewJobStore(db)
	tower := module.NewTower(db)

	defer db.Close()

	for {
		// Run crons only if tower is a leader
		isLeader, err := tower.IsLeader()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while connecting to etcd`)
			continue
		}
		// Skip if tower not a leader
		if !isLeader {
			time.Sleep(10 * time.Second)
			continue
		}

		crons, _ := cronStore.GetCronsToRun()

		// Skip if no crons
		if len(crons) == 0 {
			time.Sleep(10 * time.Second)
			continue
		}

		cron := *crons[0]

		uuid := util.GenerateUUID4()

		err = jobStore.CreateRecord(model.JobRecord{
			ID:       uuid,
			Hostname: cron.Hostname,
			CronID:   cron.ID,
			Status:   model.PendingStatus,
		})

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while connecting to etcd`)
			continue
		}

		job, err := jobStore.GetRecord(cron.Hostname, uuid)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while connecting to etcd`)
			continue
		}

		// Update last run
		cron.LastRun = time.Now().Unix()
		cronStore.UpdateRecord(cron)

		message := module.BackupMessage{
			Action:        "backup",
			Job:           *job,
			Cron:          cron,
			Settings:      map[string]string{},
			CorrelationID: util.GenerateUUID4(),
		}

		log.WithFields(log.Fields{
			"message": message,
		}).Info(`Broadcast request`)

		w.broadcast <- message

		time.Sleep(10 * time.Second)
	}
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

	db := driver.NewEtcdDriver()
	err := db.Connect()

	if err != nil {
		panic(fmt.Sprintf(
			"Error while connecting to etcd: %s",
			err.Error(),
		))
	}

	defer db.Close()

	httpClient := service.NewHTTPClient(30)
	wire := module.NewWire(httpClient, db)

	for message := range w.broadcast {

		log.WithFields(log.Fields{
			"message": message,
		}).Info(`Tower worker received a new message`)

		err = wire.SendJobToHostAgent(message)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while sending backup request`)
			continue
		}

		notifyChannel <- message
	}

	wg.Done()
}

// NotifyTower notifies tower
func (w *Workers) NotifyTower(notifyChannel <-chan module.BackupMessage) {

	for message := range notifyChannel {
		log.WithFields(log.Fields{
			"message": message,
		}).Info(`Tower worker finished processing a message`)
	}
}
