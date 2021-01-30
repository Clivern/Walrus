// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package agent

import (
	"time"

	"github.com/clivern/walrus/core/module"
	"github.com/clivern/walrus/core/service"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Heartbeat function
func Heartbeat() {
	httpClient := service.NewHTTPClient(30)
	agent := module.NewAgent(httpClient)

	for {
		log.Info(`Agent heartbeat`)

		err := agent.Heartbeat()

		if err != nil {
			log.WithFields(log.Fields{
				"error": err.Error(),
			}).Error(`Error while calling tower`)

			time.Sleep(time.Duration(viper.GetInt("agent.tower.pingInterval")) * time.Second)
			continue
		}

		time.Sleep(time.Duration(viper.GetInt("agent.tower.pingInterval")) * time.Second)
	}
}
