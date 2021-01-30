// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Host type
type Host struct {
	db driver.Database
}

// HostData type
type HostData struct {
	ID        string `json:"id"`
	Hostname  string `json:"hostname"`
	CreatedAt int64  `json:"createdAt"`
	UpdatedAt int64  `json:"updatedAt"`
}

// NewHostStore creates a new instance
func NewHostStore(db driver.Database) *Host {
	result := new(Host)
	result.db = db

	return result
}

// CreateHost creates a host
func (h *Host) CreateHost(hostData HostData) error {

	hostData.ID = util.GenerateUUID4()
	hostData.CreatedAt = time.Now().Unix()
	hostData.UpdatedAt = time.Now().Unix()

	result, err := util.ConvertToJSON(hostData)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"host_id":  hostData.ID,
		"hostname": hostData.Hostname,
	}).Debug("Create an agent")

	// store agent data
	err = h.db.Put(fmt.Sprintf(
		"%s/host/%s/h-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostData.Hostname,
	), result)

	if err != nil {
		return err
	}

	return nil
}

// UpdateHost updates a host
func (h *Host) UpdateHost(hostData HostData) error {
	hostData.UpdatedAt = time.Now().Unix()

	result, err := util.ConvertToJSON(hostData)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"host_id":  hostData.ID,
		"hostname": hostData.Hostname,
	}).Debug("Update host")

	// store agent status
	err = h.db.Put(fmt.Sprintf(
		"%s/host/%s/h-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostData.Hostname,
	), result)

	if err != nil {
		return err
	}

	return nil
}

// GetHost gets a host
func (h *Host) GetHost(hostname string) (*HostData, error) {
	hostData := &HostData{}

	log.WithFields(log.Fields{
		"hostname": hostname,
	}).Debug("Get a host")

	data, err := h.db.Get(fmt.Sprintf(
		"%s/host/%s/h-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
	))

	if err != nil {
		return hostData, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/h-data") {
			err = util.LoadFromJSON(hostData, []byte(v))

			if err != nil {
				return hostData, err
			}

			return hostData, nil
		}
	}

	return hostData, fmt.Errorf(
		"Unable to find host with name: %s",
		hostname,
	)
}

// GetHosts get hosts
func (h *Host) GetHosts() ([]*HostData, error) {

	log.Debug("Get hosts")

	records := make([]*HostData, 0)

	data, err := h.db.Get(fmt.Sprintf(
		"%s/host",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
	))

	if err != nil {
		return records, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/h-data") {
			recordData := &HostData{}

			err = util.LoadFromJSON(recordData, []byte(v))

			if err != nil {
				return records, err
			}

			records = append(records, recordData)
		}
	}

	return records, nil
}

// DeleteHost deletes a host
func (h *Host) DeleteHost(hostname string) (bool, error) {

	log.WithFields(log.Fields{
		"hostname": hostname,
	}).Debug("Delete a host")

	count, err := h.db.Delete(fmt.Sprintf(
		"%s/host/%s",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
	))

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
