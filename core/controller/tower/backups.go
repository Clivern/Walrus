// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetHostBackups controller
func GetHostBackups(c *gin.Context) {
	c.Status(http.StatusAccepted)
	return
}

// GetBackup controller
func GetBackup(c *gin.Context) {
	c.Status(http.StatusAccepted)
	return
}

// DeleteBackup controller
func DeleteBackup(c *gin.Context) {
	c.Status(http.StatusAccepted)
	return
}
