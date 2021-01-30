// Copyright 2018 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"fmt"

	"github.com/clivern/walrus/core/driver"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Stats type
type Stats struct {
	db driver.Database
}

// NewStats creates a stats instance
func NewStats(db driver.Database) *Stats {
	result := new(Stats)
	result.db = db

	return result
}

// GetTotalTowers gets total towers count
func (s *Stats) GetTotalTowers() (int, error) {

	log.Debug("Counting towers")

	key := fmt.Sprintf(
		"%s/tower",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
	)

	keys, err := s.db.GetKeys(key)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while getting towers count")

		return 0, err
	}

	log.WithFields(log.Fields{
		"count": len(keys),
	}).Debug("Current towers count")

	return len(keys), nil
}

// GetTotalHosts gets total hosts count
func (s *Stats) GetTotalHosts() (int, error) {

	log.Debug("Counting hosts")

	key := fmt.Sprintf(
		"%s/host",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
	)

	keys, err := s.db.GetKeys(key)

	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("Error while getting hosts count")

		return 0, err
	}

	log.WithFields(log.Fields{
		"count": len(keys),
	}).Debug("Current hosts count")

	return len(keys), nil
}
