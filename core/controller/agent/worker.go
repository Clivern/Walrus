// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package agent

import (
	"github.com/clivern/walrus/core/util"

	log "github.com/sirupsen/logrus"
)

// Worker controller
func Worker(workerID int, messages <-chan string) {
	log.WithFields(log.Fields{
		"correlation_id": util.GenerateUUID4(),
		"worker_id":      workerID,
	}).Info(`Worker started`)

	for message := range messages {
		log.WithFields(log.Fields{
			"worker_id": workerID,
			"message":   message,
		}).Info(`Worker received a new message`)
	}
}
