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
	"github.com/clivern/walrus/core/util"

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
		viper.SetDefault("app.name", util.GenerateUUID4())

		if viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))) != "stdout" {
			dir, _ := filepath.Split(viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))))

			if !util.DirExists(dir) {
				if _, err := util.EnsureDir(dir, 775); err != nil {
					panic(fmt.Sprintf(
						"Directory [%s] creation failed with error: %s",
						dir,
						err.Error(),
					))
				}
			}

			if !util.FileExists(viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role")))) {
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

		if viper.GetString(fmt.Sprintf("%s.mode", viper.GetString("role"))) == "prod" {
			gin.SetMode(gin.ReleaseMode)
			gin.DefaultWriter = ioutil.Discard
			gin.DisableConsoleColor()
		}

		if viper.GetString(fmt.Sprintf("%s.log.format", viper.GetString("role"))) == "json" {
			log.SetFormatter(&log.JSONFormatter{})
		} else {
			log.SetFormatter(&log.TextFormatter{})
		}

		r := gin.Default()
		workers := tower.NewWorkers()

		// Allow CORS only for development
		if viper.GetString(fmt.Sprintf("%s.mode", viper.GetString("role"))) == "dev" {
			r.Use(middleware.Cors())
		}

		r.Use(middleware.Correlation())
		r.Use(middleware.Auth())

		r.Use(middleware.Decrypt())
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

		action := r.Group("/action")
		{
			action.GET("/info", tower.Info)
			action.POST("/setup", tower.Setup)
			action.POST("/auth", tower.Auth)
		}

		apiv1 := r.Group("/api/v1")
		{
			// Hosts
			apiv1.GET("/host", tower.GetHosts)
			apiv1.GET("/host/:hostname", tower.GetHost)
			apiv1.DELETE("/host/:hostname", tower.DeleteHost)

			// Host Cron
			apiv1.POST("/host/:hostname/cron", tower.CreateHostCron)
			apiv1.GET("/host/:hostname/cron", tower.GetHostCrons)
			apiv1.GET("/host/:hostname/cron/:cronId", tower.GetHostCron)
			apiv1.PUT("/host/:hostname/cron/:cronId", tower.UpdateHostCron)
			apiv1.DELETE("/host/:hostname/cron/:cronId", tower.DeleteHostCron)

			// Host Jobs
			apiv1.GET("/host/:hostname/job", tower.GetHostJobs)
			apiv1.GET("/host/:hostname/job/:jobId", tower.GetHostJob)
			apiv1.PUT("/host/:hostname/job/:jobId", tower.UpdateHostJob)
			apiv1.DELETE("/host/:hostname/job/:jobId", tower.DeleteHostJob)

			// Settings
			apiv1.GET("/settings", tower.GetSettings)
			apiv1.PUT("/settings", tower.UpdateSettings)

			// These endpoints accept only encrypted data
			apiv1.POST("/action/bootstrap_agent", tower.AgentBootstrap)
			apiv1.POST("/action/agent_heartbeat", tower.AgentHeartbeat)
			apiv1.POST("/action/agent_postback", tower.AgentPostback)
		}

		go workers.BroadcastRequests()
		go workers.NotifyTower(workers.HandleWorkload())
		go tower.Heartbeat()

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
