// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/service"
	"github.com/clivern/walrus/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Wire type
type Wire struct {
	httpClient *service.HTTPClient
	db         driver.Database
	job        *model.Job
	agent      *model.Agent
	option     *model.Option
}

// BackupMessage type
type BackupMessage struct {
	Action        string            `json:"action"`
	Cron          model.CronRecord  `json:"cron"`
	Job           model.JobRecord   `json:"job"`
	Settings      map[string]string `json:"settings"`
	CorrelationID string            `json:"CorrelationID"`
}

// PostbackRequest type
type PostbackRequest struct {
	JobID    string `json:"jobId"`
	CronID   string `json:"cronId"`
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
	AgentID  string `json:"agentID"`
}

// NewWire creates a new instance
func NewWire(httpClient *service.HTTPClient, db driver.Database) *Wire {
	result := new(Wire)
	result.db = db
	result.httpClient = httpClient
	result.job = model.NewJobStore(db)
	result.agent = model.NewAgentStore(db)
	result.option = model.NewOptionStore(db)

	return result
}

// AgentPostback trigger agent postback. It updates job status inside a tower
func (w *Wire) AgentPostback(jobID, cronID, status string) error {

	log.Debug("Agent postback")

	hostname, err := util.GetHostname()

	if err != nil {
		return fmt.Errorf("Error while getting the hostname")
	}

	url := fmt.Sprintf(
		"%s/api/v1/action/agent_postback",
		viper.GetString("agent.tower.url"),
	)

	body, _ := util.ConvertToJSON(PostbackRequest{
		JobID:    jobID,
		CronID:   cronID,
		Status:   status,
		Hostname: hostname,
		AgentID:  viper.GetString("app.name"),
	})

	if viper.GetString("agent.tower.encryptionKey") == "" {
		return fmt.Errorf("Config agent.tower.encryptionKey is missing")
	}

	bodyByte, err := util.Encrypt(
		[]byte(body),
		viper.GetString("agent.tower.encryptionKey"),
	)

	if err != nil {
		return err
	}

	response, err := w.httpClient.Post(
		context.TODO(),
		url,
		string(bodyByte),
		map[string]string{},
		map[string]string{"X-API-Key": viper.GetString("agent.tower.apiKey"), "X-Encrypted-Request": "true"},
	)

	if err != nil {
		return err
	}

	if w.httpClient.GetStatusCode(response) != http.StatusOK {
		return fmt.Errorf(
			"Invalid response code: %d",
			w.httpClient.GetStatusCode(response),
		)
	}

	return nil
}

// UpdateTowerJobStatus updates a job status
func (w *Wire) UpdateTowerJobStatus(hostname, jobID, status string) error {

	log.Debug("Update tower job status")

	record, err := w.job.GetRecord(hostname, jobID)

	if err != nil {
		return err
	}

	record.Status = status

	err = w.job.UpdateRecord(*record)

	if err != nil {
		return err
	}

	return nil
}

// SendJobToHostAgent updates a job status
func (w *Wire) SendJobToHostAgent(message BackupMessage) error {

	agents, err := w.agent.GetAgents(message.Cron.Hostname)
	agent := &model.AgentData{}

	// TODO: Select a random running agents
	for _, v := range agents {
		if v.Status == model.UpStatus {
			agent = v
			break
		}
	}

	if agent.ID == "" {
		return fmt.Errorf(
			"Unable to find running agent for host: %s",
			message.Cron.Hostname,
		)
	}

	s3Key, err := w.option.GetOptionByKey("backup_s3_key")

	if err != nil {
		return fmt.Errorf(
			"Error while getting option backup_s3_key: %s",
			err.Error(),
		)
	}

	s3Secret, err := w.option.GetOptionByKey("backup_s3_secret")

	if err != nil {
		return fmt.Errorf(
			"Error while getting option backup_s3_secret: %s",
			err.Error(),
		)
	}

	s3Endpoint, err := w.option.GetOptionByKey("backup_s3_endpoint")

	if err != nil {
		return fmt.Errorf(
			"Error while getting option backup_s3_endpoint: %s",
			err.Error(),
		)
	}

	s3Region, err := w.option.GetOptionByKey("backup_s3_region")

	if err != nil {
		return fmt.Errorf(
			"Error while getting option backup_s3_region: %s",
			err.Error(),
		)
	}

	s3Bucket, err := w.option.GetOptionByKey("backup_s3_bucket")

	if err != nil {
		return fmt.Errorf(
			"Error while getting option backup_s3_bucket: %s",
			err.Error(),
		)
	}

	message.Settings["backup_s3_key"] = s3Key.Value
	message.Settings["backup_s3_secret"] = s3Secret.Value
	message.Settings["backup_s3_endpoint"] = s3Endpoint.Value
	message.Settings["backup_s3_region"] = s3Region.Value
	message.Settings["backup_s3_bucket"] = s3Bucket.Value

	url := fmt.Sprintf(
		"%s/api/v1/process",
		agent.URL,
	)

	body, _ := util.ConvertToJSON(message)

	if viper.GetString("agent.tower.encryptionKey") == "" {
		return fmt.Errorf("Config agent.tower.encryptionKey is missing")
	}

	bodyByte, err := util.Encrypt(
		[]byte(body),
		viper.GetString("agent.tower.encryptionKey"),
	)

	if err != nil {
		return err
	}

	response, err := w.httpClient.Post(
		context.TODO(),
		url,
		string(bodyByte),
		map[string]string{},
		map[string]string{"X-API-Key": agent.APIKey, "X-Encrypted-Request": "true"},
	)

	if err != nil {
		return err
	}

	if w.httpClient.GetStatusCode(response) != http.StatusAccepted {
		return fmt.Errorf(
			"Invalid response code: %d",
			w.httpClient.GetStatusCode(response),
		)
	}

	return nil
}
