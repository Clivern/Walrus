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

	"github.com/clivern/walrus/core/controller/agent"
	"github.com/clivern/walrus/core/middleware"
	"github.com/clivern/walrus/core/service"

	"github.com/drone/envsubst"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var agentCmd = &cobra.Command{
	Use:   "agent",
	Short: "Start walrus agent",
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

		viper.SetDefault("role", "agent")

		if viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))) != "stdout" {
			dir, _ := filepath.Split(viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role"))))

			if !service.DirExists(dir) {
				if _, err := service.EnsureDir(dir, 777); err != nil {
					panic(fmt.Sprintf(
						"Directory [%s] creation failed with error: %s",
						dir,
						err.Error(),
					))
				}
			}

			if !service.FileExists(viper.GetString(fmt.Sprintf("%s.log.output", viper.GetString("role")))) {
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

		messages := make(chan string, viper.GetInt(
			fmt.Sprintf("%s.broker.native.capacity", viper.GetString("role")),
		))

		go agent.Heartbeat(messages)

		for i := 0; i < viper.GetInt(fmt.Sprintf("%s.broker.native.workers", viper.GetString("role"))); i++ {
			go agent.Worker(i+1, messages)
		}

		r := gin.Default()

		r.Use(middleware.Correlation())
		r.Use(middleware.Logger())

		r.GET("/favicon.ico", func(c *gin.Context) {
			c.String(http.StatusNoContent, "")
		})

		r.GET("/", agent.Health)

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
	agentCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	agentCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(agentCmd)
}
