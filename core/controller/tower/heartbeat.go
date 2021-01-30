// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"fmt"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/module"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	totalTowers = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "walrus",
			Name:      "cluster_total_towers",
			Help:      "Total towers in the cluster",
		})

	totalHosts = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "walrus",
			Name:      "cluster_total_hosts",
			Help:      "Total hosts in the cluster",
		})
)

func init() {
	prometheus.MustRegister(totalTowers)
	prometheus.MustRegister(totalHosts)
}

// Heartbeat node heartbeat
func Heartbeat() {
	db := driver.NewEtcdDriver()

	err := db.Connect()

	if err != nil {
		panic(fmt.Sprintf(
			"Error while connecting to etcd: %s",
			err.Error(),
		))
	}

	defer db.Close()

	stats := module.NewStats(db)
	tower := module.NewTower(db)

	log.Info(`Start heartbeat daemon`)

	count := 0

	for {
		log.Debug(`Mark the tower as a live`)

		err := tower.Alive(5)

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while connecting to etcd`)
			continue
		}

		log.Debug(`Check current cluster leader and elect a new one`)

		hasLeader, err := tower.HasLeader()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while connecting to etcd`)
			continue
		}

		isLeader, err := tower.IsLeader()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while connecting to etcd`)
			continue
		}

		// If there is no leader or the current is the leader
		// Refresh the election
		if !hasLeader || isLeader {
			err := tower.Elect(5)

			if err != nil {
				log.WithFields(log.Fields{
					"error": err.Error(),
				}).Error(`Error while connecting to etcd`)
				continue
			}
		}

		// Refresh metrics
		log.Debug(`Refresh metrics`)

		count, err = stats.GetTotalTowers()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while connecting to etcd`)
			continue
		}

		totalTowers.Set(float64(count))

		count, err = stats.GetTotalHosts()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while connecting to etcd`)
			continue
		}

		totalHosts.Set(float64(count))

		time.Sleep(2 * time.Second)

		log.Info(`Tower heartbeat done`)
	}
}
