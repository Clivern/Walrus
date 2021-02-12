// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/clivern/walrus/core/driver"
	"github.com/clivern/walrus/core/util"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	// SecondInterval constant
	SecondInterval = "@second"

	// MinuteInterval constant
	MinuteInterval = "@minute"

	// HourInterval constant
	HourInterval = "@hour"

	// DayInterval constant
	DayInterval = "@day"

	// MonthInterval constant
	MonthInterval = "@month"

	// BackupDirectory constant
	BackupDirectory = "@BackupDirectory"

	// BackupMySQL constant
	BackupMySQL = "@BackupMySQL"

	// BackupRedis constant
	// TODO Redis Backup Support
	BackupRedis = "@BackupRedis"

	// BackupPostgreSQL constant
	// TODO PostgreSQL Backup Support
	BackupPostgreSQL = "@BackupPostgreSQL"

	// BackupEtcd constant
	// TODO Etcd Backup Support
	BackupEtcd = "@BackupEtcd"

	// BackupSQLite constant
	BackupSQLite = "@BackupSQLite"
)

// Request type
type Request struct {
	Type          string `json:"type"` // @BackupSQLite, @BackupMySQL, @BackupDirectory, @BackupRedis, @BackupPostgreSQL
	Directory     string `json:"directory"`
	RetentionDays int    `json:"retentionDays"`

	MySQLHost         string `json:"mysqlHost"`
	MySQLPort         string `json:"mysqlPort"`
	MySQLUsername     string `json:"mysqlUsername"`
	MySQLPassword     string `json:"mysqlPassword"`
	MySQLAllDatabases bool   `json:"mysqlAllDatabases"`
	MySQLDatabase     string `json:"mysqlDatabase"`
	MySQLTable        string `json:"mysqlTable"`
	MySQLOptions      string `json:"mysqlOptions"`

	SQLitePath string `json:"sqlitePath"`
}

// CronRecord type
type CronRecord struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	Hostname     string  `json:"hostname"`
	Request      Request `json:"request"`
	Interval     int     `json:"interval"`
	IntervalType string  `json:"intervalType"`
	SuccessJobs  int     `json:"successJobs"`
	FailedJobs   int     `json:"failedJobs"`
	LastRun      int64   `json:"lastRun"`
	CreatedAt    int64   `json:"createdAt"`
	UpdatedAt    int64   `json:"updatedAt"`
}

// Cron type
type Cron struct {
	db driver.Database
}

// NewCronStore creates a new instance
func NewCronStore(db driver.Database) *Cron {
	result := new(Cron)
	result.db = db

	return result
}

// CreateRecord stores a cron record
func (c *Cron) CreateRecord(record CronRecord) error {
	record.CreatedAt = time.Now().Unix()
	record.UpdatedAt = time.Now().Unix()
	record.LastRun = time.Now().Unix()

	result, err := util.ConvertToJSON(record)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"cron_id":  record.ID,
		"hostname": record.Hostname,
	}).Debug("Create a record")

	// store cron record data
	err = c.db.Put(fmt.Sprintf(
		"%s/host/%s/cron/%s/c-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		record.Hostname,
		record.ID,
	), result)

	if err != nil {
		return err
	}

	return nil
}

// UpdateRecord updates a cron record
func (c *Cron) UpdateRecord(record CronRecord) error {
	record.UpdatedAt = time.Now().Unix()

	result, err := util.ConvertToJSON(record)

	if err != nil {
		return err
	}

	log.WithFields(log.Fields{
		"cron_id":  record.ID,
		"hostname": record.Hostname,
	}).Debug("Update a record")

	// store cron record data
	err = c.db.Put(fmt.Sprintf(
		"%s/host/%s/cron/%s/c-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		record.Hostname,
		record.ID,
	), result)

	if err != nil {
		return err
	}

	return nil
}

// GetRecord gets cron record data
func (c *Cron) GetRecord(hostname, cronID string) (*CronRecord, error) {
	recordData := &CronRecord{}

	log.WithFields(log.Fields{
		"cron_id":  cronID,
		"hostname": hostname,
	}).Debug("Get a record data")

	data, err := c.db.Get(fmt.Sprintf(
		"%s/host/%s/cron/%s/c-data",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
		cronID,
	))

	if err != nil {
		return recordData, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/c-data") {
			err = util.LoadFromJSON(recordData, []byte(v))

			if err != nil {
				return recordData, err
			}

			return recordData, nil
		}
	}

	return recordData, fmt.Errorf(
		"Unable to find cron record with id: %s and hostname: %s",
		cronID,
		hostname,
	)
}

// DeleteRecord deletes a cron record
func (c *Cron) DeleteRecord(hostname, cronID string) (bool, error) {

	log.WithFields(log.Fields{
		"cron_id":  cronID,
		"hostname": hostname,
	}).Debug("Delete a record")

	count, err := c.db.Delete(fmt.Sprintf(
		"%s/host/%s/cron/%s",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
		cronID,
	))

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// GetCronsToRun get crons to run
func (c *Cron) GetCronsToRun() ([]*CronRecord, error) {

	log.Debug("Get crons to run")

	records := make([]*CronRecord, 0)

	data, err := c.db.Get(fmt.Sprintf(
		"%s/host",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
	))

	if err != nil {
		return records, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/c-data") {
			recordData := &CronRecord{}

			err = util.LoadFromJSON(recordData, []byte(v))

			if err != nil {
				return records, err
			}

			tm := time.Unix(recordData.LastRun, 0)

			if recordData.IntervalType == SecondInterval {
				tm = tm.Add(time.Duration(recordData.Interval) * time.Second)
			} else if recordData.IntervalType == MinuteInterval {
				tm = tm.Add(time.Duration(recordData.Interval) * time.Minute)
			} else if recordData.IntervalType == HourInterval {
				tm = tm.Add(time.Duration(recordData.Interval) * time.Hour)
			} else if recordData.IntervalType == DayInterval {
				tm = tm.Add(time.Duration(recordData.Interval*24) * time.Hour)
			} else if recordData.IntervalType == MonthInterval {
				tm = tm.Add(time.Duration(recordData.Interval*30*24) * time.Hour)
			}

			// If time in future, skip
			if tm.After(time.Now()) {
				continue
			}

			records = append(records, recordData)
		}
	}

	return records, nil
}

// GetHostCrons get crons for a host
func (c *Cron) GetHostCrons(hostname string) ([]*CronRecord, error) {

	log.Debug("Get crons to run")

	records := make([]*CronRecord, 0)

	data, err := c.db.Get(fmt.Sprintf(
		"%s/host/%s/cron",
		viper.GetString(fmt.Sprintf("%s.database.etcd.databaseName", viper.GetString("role"))),
		hostname,
	))

	if err != nil {
		return records, err
	}

	for k, v := range data {
		// Check if it is the data key
		if strings.Contains(k, "/c-data") {
			recordData := &CronRecord{}

			err = util.LoadFromJSON(recordData, []byte(v))

			if err != nil {
				return records, err
			}

			records = append(records, recordData)
		}
	}

	return records, nil
}
