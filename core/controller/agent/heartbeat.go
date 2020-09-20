// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package agent

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// Heartbeat function
func Heartbeat(messages chan<- string) {
	for {
		log.Info(`Agent heartbeat`)

		messages <- "BAM!"

		time.Sleep(2 * time.Second)
	}
}
