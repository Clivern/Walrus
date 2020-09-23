// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cmd

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/clivern/walrus/core/controller/tower"
	"github.com/clivern/walrus/core/middleware"
	"github.com/clivern/walrus/core/service"

	"github.com/drone/envsubst"
	"github.com/gin-gonic/gin"
	"github.com/markbates/pkger"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var towerCmd = &cobra.Command{
	Use:   "tower",
	Short: "Start walrus tower",
	Run: func(cmd *cobra.Command, args []string) {
		configUnparsed, err := ioutil.ReadFile(config)

		if err != nil {
			panic(fmt.Sprintf(
				"Error while reading config file [%s]: %s",
				config,
				err.Error(),
			))
		}

		configParsed, err := envsubst.EvalEnv(string(configUnparsed))

		if err != nil {
			panic(fmt.Sprintf(
				"Error while parsing config file [%s]: %s",
				config,
				err.Error(),
			))
		}

		viper.SetConfigType("yaml")
		err = viper.ReadConfig(bytes.NewBuffer([]byte(configParsed)))

		if err != nil {
			panic(fmt.Sprintf(
				"Error while loading configs [%s]: %s",
				config,
				err.Error(),
			))
		}

		viper.SetDefault("role", "tower")

		if viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))) != "stdout" {
			fs := service.FileSystem{}
			dir, _ := filepath.Split(viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))))

			if !fs.DirExists(dir) {
				if _, err := fs.EnsureDir(dir, 777); err != nil {
					panic(fmt.Sprintf(
						"Directory [%s] creation failed with error: %s",
						dir,
						err.Error(),
					))
				}
			}

			if !fs.FileExists(viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role")))) {
				f, err := os.Create(viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))))
				if err != nil {
					panic(fmt.Sprintf(
						"Error while creating log file [%s]: %s",
						viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))),
						err.Error(),
					))
				}
				defer f.Close()
			}
		}

		if viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))) == "stdout" {
			gin.DefaultWriter = os.Stdout
			log.SetOutput(os.Stdout)
		} else {
			f, _ := os.Create(viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))))
			gin.DefaultWriter = io.MultiWriter(f)
			log.SetOutput(f)
		}

		lvl := strings.ToLower(viper.GetString(fmt.Sprintf("%s.log.level", viper.GetString("role"))))
		level, err := log.ParseLevel(lvl)

		if err != nil {
			level = log.InfoLevel
		}

		log.SetLevel(level)

		if viper.GetString(fmt.Sprintf("%s.app.mode", viper.GetString("role"))) == "prod" {
			gin.SetMode(gin.ReleaseMode)
			gin.DefaultWriter = ioutil.Discard
			gin.DisableConsoleColor()
		}

		if viper.GetString(fmt.Sprintf("%s.log.format", viper.GetString("role"))) == "json" {
			log.SetFormatter(&log.JSONFormatter{})
		} else {
			log.SetFormatter(&log.TextFormatter{})
		}

		// Init DB Connection
		db := service.Database{}
		err = db.AutoConnect()

		if err != nil {
			panic(err.Error())
		}

		// Migrate Database
		success := db.Migrate()

		if !success {
			panic("Error! Unable to migrate database tables.")
		}

		defer db.Close()

		r := gin.Default()

		r.Use(middleware.Correlation())
		r.Use(middleware.Auth())

		// Allow CORS only for development
		if viper.GetString(fmt.Sprintf("%s.app.mode", viper.GetString("role"))) == "dev" {
			r.Use(middleware.Cors())
		}

		r.Use(middleware.Logger())
		r.Use(middleware.Metric())

		r.GET("/favicon.ico", func(c *gin.Context) {
			c.String(http.StatusNoContent, "")
		})

		r.GET("/", tower.Home)
		r.GET("/_health", tower.Health)
		r.GET("/_ready", tower.Ready)

		r.GET(
			viper.GetString(fmt.Sprintf("%s.metrics.prometheus.endpoint", viper.GetString("role"))),
			gin.WrapH(tower.Metrics()),
		)

		r.NoRoute(gin.WrapH(http.FileServer(pkger.Dir("/web/dist"))))

		r.GET("/info", tower.Info)
		r.POST("/setup", tower.Setup)
		r.POST("/auth", tower.Auth)

		r.GET("/api/v1/jobs", tower.GetJobs)
		r.GET("/api/v1/jobs/:jobId", tower.GetJob)
		r.DELETE("/api/v1/jobs/:jobId", tower.DeleteJob)

		// These endpoints accept only encrypted data
		r.POST("/api/v1/agent/heartbeat", tower.AgentsHeartbeat)
		r.POST("/api/v1/agent/postback", tower.AgentsPostback)

		r.GET("/api/v1/hosts", tower.GetHosts)
		r.GET("/api/v1/hosts/:hostId", tower.GetHost)
		r.PUT("/api/v1/hosts/:hostId", tower.UpdateHost)
		r.DELETE("/api/v1/hosts/:hostId", tower.DeleteHost)

		r.GET("/api/v1/hosts/:hostId/backups", tower.GetHostBackups)
		r.GET("/api/v1/hosts/:hostId/backups/:backupId", tower.GetBackup)
		r.DELETE("/api/v1/hosts/:hostId/backups/:backupId", tower.DeleteBackup)

		r.GET("/api/v1/settings", tower.GetSettings)
		r.PUT("/api/v1/settings", tower.UpdateSettings)

		var runerr error

		if viper.GetBool(fmt.Sprintf("%s.tls.status", viper.GetString("role"))) {
			runerr = r.RunTLS(
				fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt(fmt.Sprintf("%s.port", viper.GetString("role"))))),
				viper.GetString(fmt.Sprintf("%s.tls.pemPath", viper.GetString("role"))),
				viper.GetString(fmt.Sprintf("%s.tls.keyPath", viper.GetString("role"))),
			)
		} else {
			runerr = r.Run(
				fmt.Sprintf(":%s", strconv.Itoa(viper.GetInt(fmt.Sprintf("%s.port", viper.GetString("role"))))),
			)
		}

		if runerr != nil {
			panic(runerr.Error())
		}
	},
}

func init() {
	towerCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	towerCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(towerCmd)
}
