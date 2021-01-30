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

const (
	// UnknownStatus const
	UnknownStatus = "unknown"
	// UpStatus const
	UpStatus = "up"
	// DownStatus const
	DownStatus = "down"
)

// Agent type
type Agent struct {
	db driver.Database
}

// AgentData type
type AgentData struct {
	ID              string `json:"id"`
	URL             string `json:"url"`
	Hostname        string `json:"hostname"`
	APIKey          string `json:"apiKey"`
	Status          string `json:"status"`
	CreatedAt       int64  `json:"createdAt"`
	UpdatedAt       int64  `json:"updatedAt"`
	LastStatusCheck int64  `json:"lastStatusCheck"`
}

// NewAgentStore creates a new instance
func NewAgentStore(db driver.Database) *Agent {
	result := new(Agent)
	result.db = db

	return result
}

// CreateAgent stores agent data and status
func (a *Agent) CreateAgent(agentData AgentData) error {

	agentData.CreatedAt = time.Now().Unix()
	agentData.UpdatedAt = time.Now().Unix()
	agentData.LastStatusCheck = time.Now().Unix()
	agentData.Status = UnknownStatus

	result, err := util.ConvertToJSON(agentData)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"agent_id": agentData.ID,
		"hostname": agentData.Hostname,
	}).Debug("Create an agent")

	// store agent data
	err = a.db.Put(fmt.Sprintf(
		"%s/host/%s/agent/%s/a-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		agentData.Hostname,
		agentData.ID,
	), result)

	if err != nil {
		return err
	}

	return nil
}

// UpdateAgent updates agent data
func (a *Agent) UpdateAgent(agentData AgentData) error {
	agentData.UpdatedAt = time.Now().Unix()

	result, err := util.ConvertToJSON(agentData)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"agent_id": agentData.ID,
		"hostname": agentData.Hostname,
	}).Debug("Update agent data")

	// update agent data
	err = a.db.Put(fmt.Sprintf(
		"%s/host/%s/agent/%s/a-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		agentData.Hostname,
		agentData.ID,
	), result)

	if err != nil {
		return err
	}

	return nil
}

// GetAgent gets agent data
func (a *Agent) GetAgent(hostname, agentID string) (*AgentData, error) {
	agentData := &AgentData{}

	log.WithFields(log.Fields{
		"agent_id": agentID,
		"hostname": hostname,
	}).Debug("Get agent data")

	data, err := a.db.Get(fmt.Sprintf(
		"%s/host/%s/agent/%s/a-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
		agentID,
	))

	if err != nil {
		return agentData, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/a-data") {
			err = util.LoadFromJSON(agentData, []byte(v))

			if err != nil {
				return agentData, err
			}

			return agentData, nil
		}
	}

	return agentData, fmt.Errorf(
		"Unable to find agent with id: %s and hostname: %s",
		agentID,
		hostname,
	)
}

// GetAgents get agents
func (a *Agent) GetAgents(hostname string) ([]*AgentData, error) {

	records := make([]*AgentData, 0)

	log.WithFields(log.Fields{
		"hostname": hostname,
	}).Debug("Get agents")

	data, err := a.db.Get(fmt.Sprintf(
		"%s/host/%s",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
	))

	if err != nil {
		return records, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/a-data") {
			recordData := &AgentData{}

			err = util.LoadFromJSON(recordData, []byte(v))

			if err != nil {
				return records, err
			}

			records = append(records, recordData)
		}
	}

	return records, nil
}

// CountOnlineAgents counts online agents
func (a *Agent) CountOnlineAgents(hostname string) (int, error) {

	count := 0

	log.WithFields(log.Fields{
		"hostname": hostname,
	}).Debug("Count online agents")

	data, err := a.db.Get(fmt.Sprintf(
		"%s/host/%s",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
	))

	if err != nil {
		return count, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/a-data") {
			recordData := &AgentData{}

			err = util.LoadFromJSON(recordData, []byte(v))

			if err != nil {
				return count, err
			}

			// Agent become inactive if heartbeat missing for more than
			// 10 min
			if recordData.Status == UpStatus && recordData.LastStatusCheck >= time.Now().Add(-10*time.Minute).Unix() {
				count++
			}
		}
	}

	return count, nil
}

// DeleteAgent deletes an agent by id and hostname
func (a *Agent) DeleteAgent(hostname, agentID string) (bool, error) {

	log.WithFields(log.Fields{
		"agent_id": agentID,
		"hostname": hostname,
	}).Debug("Delete an agent")

	count, err := a.db.Delete(fmt.Sprintf(
		"%s/host/%s/agent/%s",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
		agentID,
	))

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
