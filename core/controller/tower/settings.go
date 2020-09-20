// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetSettings controller
func GetSettings(c *gin.Context) {
	c.Status(http.StatusAccepted)
	return
}

// UpdateSettings controller
func UpdateSettings(c *gin.Context) {
	c.Status(http.StatusAccepted)
	return
}
