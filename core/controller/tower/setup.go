// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package tower

import (
	"io/ioutil"
	"net/http"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/model"
	"github.com/clivern/walrus/core/util"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Setup controller
func Setup(c *gin.Context) {
	var userData model.UserData

	data, _ := ioutil.ReadAll(c.Request.Body)

	err := util.LoadFromJSON(&userData, data)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Error! Invalid request",
		})
		return
	}

	if !util.IsEmailValid(userData.Email) || util.IsEmpty(userData.Password) || util.IsEmpty(userData.Name) {
		c.JSON(http.StatusBadRequest, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Please provide all required data",
		})
		return
	}

	db := driver.NewEtcdDriver()

	err = db.Connect()

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	defer db.Close()

	userStore := model.NewUserStore(db)

	optionStore := model.NewOptionStore(db)

	_, err = optionStore.GetOptionByKey("is_installed")

	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Error! Walrus installed already",
		})
		return
	}

	err = userStore.CreateUser(model.UserData{
		ID:       util.GenerateUUID4(),
		Name:     userData.Name,
		Email:    userData.Email,
		Password: userData.Password,
		APIKey:   util.GenerateUUID4(),
	})

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	err = optionStore.CreateOption(model.OptionData{
		Key:   "is_installed",
		Value: "true",
	})

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	err = optionStore.CreateOption(model.OptionData{
		Key:   "backup_s3_key",
		Value: "",
	})

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	err = optionStore.CreateOption(model.OptionData{
		Key:   "backup_s3_secret",
		Value: "",
	})

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	err = optionStore.CreateOption(model.OptionData{
		Key:   "backup_s3_endpoint",
		Value: "https://nyc3.digitaloceanspaces.com",
	})

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	err = optionStore.CreateOption(model.OptionData{
		Key:   "backup_s3_region",
		Value: "us-east-1",
	})

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	err = optionStore.CreateOption(model.OptionData{
		Key:   "backup_s3_bucket",
		Value: util.GenerateUUID4(),
	})

	if err != nil {
		log.WithFields(log.Fields{
			"correlation_id": c.GetHeader("x-correlation-id"),
			"error":          err.Error(),
		}).Error("Internal server error")

		c.JSON(http.StatusInternalServerError, gin.H{
			"correlationID": c.GetHeader("x-correlation-id"),
			"errorMessage":  "Internal server error",
		})
		return
	}

	c.Status(http.StatusOK)
}
