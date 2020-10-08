// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"io/ioutil"
	"net/http"

	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/service"
	"github.com/clivern/walrus/core/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Setup controller
func Setup(c *gin.Context) {
	request := make(map[string]string)

	db := service.Database{}

	data, _ := ioutil.ReadAll(c.Request.Body)

	util.LoadFromJSON(&request, data)

	err := db.AutoConnect()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.Request.Header.Get("X-Correlation-ID"),
			"error":          err.Error(),
		}).Error(`Failure while connecting database`)

		c.Status(http.StatusInternalServerError)
		return
	}

	defer db.Close()

	request["email"]
	request["password"]


	c.Status(http.StatusOK)
	return
}
