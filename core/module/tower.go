// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"fmt"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Tower type
type Tower struct {
	db driver.Database
}

// NewTower creates a new instance
func NewTower(db driver.Database) *Tower {
	result := new(Tower)
	result.db = db

	return result
}

// Elect elect as leader
func (t *Tower) Elect(seconds int64) error {

	log.Debug("Elect a leader")

	hostname, err := util.GetHostname()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while getting the hostname")

		return err
	}

	key := fmt.Sprintf(
		"%s/leader/id",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
	)

	value := fmt.Sprintf(
		"%s__%s",
		hostname,
		viper.GetString("app.name"),
	)

	leaseID, err := t.db.CreateLease(seconds)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while creating a lease")

		return err
	}

	err = t.db.PutWithLease(key, value, leaseID)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while creating a key with lease")

		return err
	}

	err = t.db.RenewLease(leaseID)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while renewing a lease")

		return err
	}

	return nil
}

// IsLeader checks if current tower is a leader
func (t *Tower) IsLeader() (bool, error) {

	log.Debug("Check if tower is a leader")

	hostname, err := util.GetHostname()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while getting the hostname")

		return false, err
	}

	key := fmt.Sprintf(
		"%s/leader/id",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
	)

	value := fmt.Sprintf(
		"%s__%s",
		hostname,
		viper.GetString("app.name"),
	)

	data, err := t.db.Get(key)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"key":   key,
		}).Error("Error while getting value for a key")

		return false, err
	}

	if val, ok := data[key]; ok {
		// Is leader
		if val == value {
			log.WithFields(log.Fields{
				"tower_id": value,
			}).Debug("Tower is a leader")

			return true, nil
		}

		log.WithFields(log.Fields{
			"tower_id":  value,
			"leader_id": val,
		}).Debug("Tower is not a leader")
	}

	return false, nil
}

// HasLeader checks if there is a leader
func (t *Tower) HasLeader() (bool, error) {

	log.Debug("Check if there is a leader")

	key := fmt.Sprintf(
		"%s/leader/id",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
	)

	data, err := t.db.Get(key)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
			"key":   key,
		}).Error("Error while getting value for a key")

		return false, err
	}

	if val, ok := data[key]; ok {

		log.WithFields(log.Fields{
			"leader_id": val,
		}).Debug("Cluster has a leader")

		return true, nil
	}

	return false, nil
}

// Alive report the tower as live to etcd
func (t *Tower) Alive(seconds int64) error {

	log.Debug("Mark tower as a live")

	hostname, err := util.GetHostname()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while getting the hostname")

		return err
	}

	key := fmt.Sprintf(
		"%s/tower/%s__%s",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
		viper.GetString("app.name"),
	)

	leaseID, err := t.db.CreateLease(seconds)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while creating a lease")

		return err
	}

	err = t.db.PutWithLease(fmt.Sprintf("%s/state", key), "alive", leaseID)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while creating a key with lease")

		return err
	}

	err = t.db.RenewLease(leaseID)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while renewing a lease")

		return err
	}

	return nil
}
