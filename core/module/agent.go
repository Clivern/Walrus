// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package module

import (
	"context"
	"fmt"
	"net/http"

	"github.com/clivern/walrus/core/service"
	"github.com/clivern/walrus/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Agent type
type Agent struct {
	httpClient *service.HTTPClient
}

// BootstrapRequest type
type BootstrapRequest struct {
	AgentURL    string `json:"agentURL"`
	Hostname    string `json:"hostname"`
	AgentID     string `json:"agentID"`
	AgentAPIKey string `json:"agentApiKey"`
}

// HeartbeatRequest type
type HeartbeatRequest struct {
	Status   string `json:"status"`
	Hostname string `json:"hostname"`
	AgentID  string `json:"agentID"`
}

// NewAgent creates a new instance
func NewAgent(httpClient *service.HTTPClient) *Agent {
	return &Agent{
		httpClient: httpClient,
	}
}

// Bootstrap registers an agent
// In order to register an agent, agent do a request to tower with the following specs
//
// POST: config(agent.tower.url)/api/v1/action/register_agent
// {
//     "agentURL": "config(agent.url)"
//     "agentID": "$$",
//     "hostname": "$$"
//     "agentApiKey": "$$"
// }
//
// Headers
// X-Encrypted-Request: true
// X-API-Key: config(agent.tower.apiKey)
func (a *Agent) Bootstrap() error {

	log.Debug("Bootstrap agent")

	hostname, err := util.GetHostname()

	if err != nil {
		return fmt.Errorf("Error while getting the hostname")
	}

	url := fmt.Sprintf(
		"%s/api/v1/action/bootstrap_agent",
		viper.GetString("agent.tower.url"),
	)

	if viper.GetString("agent.tower.encryptionKey") == "" {
		return fmt.Errorf("Config agent.tower.encryptionKey is missing")
	}

	body, _ := util.ConvertToJSON(BootstrapRequest{
		AgentURL:    viper.GetString("agent.url"),
		Hostname:    hostname,
		AgentID:     viper.GetString("app.name"),
		AgentAPIKey: viper.GetString("agent.api.key"),
	})

	bodyByte, err := util.Encrypt(
		[]byte(body),
		viper.GetString("agent.tower.encryptionKey"),
	)

	if err != nil {
		return err
	}

	response, err := a.httpClient.Post(
		context.TODO(),
		url,
		string(bodyByte),
		map[string]string{},
		map[string]string{"X-API-Key": viper.GetString("agent.tower.apiKey"), "X-Encrypted-Request": "true"},
	)

	if err != nil {
		return err
	}

	if a.httpClient.GetStatusCode(response) != http.StatusOK {
		return fmt.Errorf(
			"Invalid response code: %d",
			a.httpClient.GetStatusCode(response),
		)
	}

	return nil
}

// Heartbeat notify the tower that agent is a live
//
// POST: config(agent.tower.url)/api/v1/action/agent_heartbeat
// {
//     "status": "up",
//     "agentID": "$$",
//     "hostname": "$$"
// }
//
// Headers
// X-Encrypted-Request: true
// X-API-Key: config(agent.tower.apiKey)
func (a *Agent) Heartbeat() error {

	log.Debug("Agent heartbeat")

	hostname, err := util.GetHostname()

	if err != nil {
		return fmt.Errorf("Error while getting the hostname")
	}

	url := fmt.Sprintf(
		"%s/api/v1/action/agent_heartbeat",
		viper.GetString("agent.tower.url"),
	)

	if viper.GetString("agent.tower.encryptionKey") == "" {
		return fmt.Errorf("Config agent.tower.encryptionKey is missing")
	}

	body, _ := util.ConvertToJSON(HeartbeatRequest{
		Status:   "up",
		Hostname: hostname,
		AgentID:  viper.GetString("app.name"),
	})

	bodyByte, err := util.Encrypt(
		[]byte(body),
		viper.GetString("agent.tower.encryptionKey"),
	)

	if err != nil {
		return err
	}

	response, err := a.httpClient.Post(
		context.TODO(),
		url,
		string(bodyByte),
		map[string]string{},
		map[string]string{"X-API-Key": viper.GetString("agent.tower.apiKey"), "X-Encrypted-Request": "true"},
	)

	if err != nil {
		return err
	}

	if a.httpClient.GetStatusCode(response) != http.StatusOK {
		return fmt.Errorf(
			"Invalid response code: %d",
			a.httpClient.GetStatusCode(response),
		)
	}

	return nil
}
