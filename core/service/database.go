// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package service

import (
	"fmt"
	"github.com/clivern/walrus/core/migration"
	"github.com/clivern/walrus/core/model"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	log "github.com/sirupsen/logrus"
)

// Database struct
type Database struct {
	Connection *gorm.DB
}

// Connect connects to a MySQL database
func (db *Database) Connect(dsn model.DSN) error {
	var err error

	// Reuse db connections http://go-database-sql.org/surprises.html
	if db.Ping() == nil {
		return nil
	}

	db.Connection, err = gorm.Open(dsn.Driver, dsn.ToString())

	if err != nil {
		return err
	}

	return nil
}

// Ping check the db connection
func (db *Database) Ping() error {

	if db.Connection == nil {
		return fmt.Errorf("No DB Connections Found")
	}

	err := db.Connection.DB().Ping()

	if err != nil {
		return err
	}

	// Cleanup stale connections http://go-database-sql.org/surprises.html
	db.Connection.DB().SetMaxOpenConns(5)
	db.Connection.DB().SetConnMaxLifetime(time.Duration(10) * time.Second)
	dbStats := db.Connection.DB().Stats()

	log.WithFields(log.Fields{
		"dbStats.maxOpenConnections": int(dbStats.MaxOpenConnections),
		"dbStats.openConnections":    int(dbStats.OpenConnections),
		"dbStats.inUse":              int(dbStats.InUse),
		"dbStats.idle":               int(dbStats.Idle),
	}).Debug(`Open DB Connection`)

	return nil
}

// AutoConnect connects to a MySQL database using loaded configs
func (db *Database) AutoConnect() error {
	var err error

	// Reuse db connections http://go-database-sql.org/surprises.html
	if db.Ping() == nil {
		return nil
	}

	dsn := model.DSN{
		Driver:   viper.GetString(fmt.Sprintf("%s.database.driver", viper.GetString("role"))),
		Username: viper.GetString(fmt.Sprintf("%s.database.username", viper.GetString("role"))),
		Password: viper.GetString(fmt.Sprintf("%s.database.password", viper.GetString("role"))),
		Hostname: viper.GetString(fmt.Sprintf("%s.database.host", viper.GetString("role"))),
		Port:     viper.GetInt(fmt.Sprintf("%s.database.port", viper.GetString("role"))),
		Name:     viper.GetString(fmt.Sprintf("%s.database.name", viper.GetString("role"))),
	}

	db.Connection, err = gorm.Open(dsn.Driver, dsn.ToString())

	if err != nil {
		return err
	}

	return nil
}

// Migrate migrates the database
func (db *Database) Migrate() bool {
	status := true
	db.Connection.AutoMigrate(
		&migration.Job{},
		&migration.Host{},
		&migration.Option{},
		&migration.HostMeta{},
	)

	status = status && db.Connection.HasTable(&migration.Job{})
	status = status && db.Connection.HasTable(&migration.Host{})
	status = status && db.Connection.HasTable(&migration.Option{})
	status = status && db.Connection.HasTable(&migration.HostMeta{})
	return status
}

// Rollback drop tables
func (db *Database) Rollback() bool {
	status := true
	db.Connection.DropTableIfExists(
		&migration.Job{},
		&migration.Host{},
		&migration.Option{},
		&migration.HostMeta{},
	)

	status = status && !db.Connection.HasTable(&migration.Job{})
	status = status && !db.Connection.HasTable(&migration.Host{})
	status = status && !db.Connection.HasTable(&migration.Option{})
	status = status && !db.Connection.HasTable(&migration.HostMeta{})
	return status
}

// HasTable checks if table exists
func (db *Database) HasTable(table string) bool {
	return db.Connection.HasTable(table)
}

// CreateJob creates a new job
func (db *Database) CreateJob(job *model.Job) *model.Job {
	db.Connection.Create(job)
	return job
}

// JobExistByID check if job exists
func (db *Database) JobExistByID(id int) bool {
	job := model.Job{}

	db.Connection.Where("id = ?", id).First(&job)

	return job.ID > 0
}

// GetJobByID gets a job by id
func (db *Database) GetJobByID(id int) model.Job {
	job := model.Job{}

	db.Connection.Where("id = ?", id).First(&job)

	return job
}

// GetJobs gets jobs
func (db *Database) GetJobs() []model.Job {
	jobs := []model.Job{}

	db.Connection.Select("*").Find(&jobs)

	return jobs
}

// JobExistByUUID check if job exists
func (db *Database) JobExistByUUID(uuid string) bool {
	job := model.Job{}

	db.Connection.Where("uuid = ?", uuid).First(&job)

	return job.ID > 0
}

// GetJobByUUID gets a job by uuid
func (db *Database) GetJobByUUID(uuid string) model.Job {
	job := model.Job{}

	db.Connection.Where("uuid = ?", uuid).First(&job)

	return job
}

// GetPendingJobByType gets a job by uuid
func (db *Database) GetPendingJobByType(jobType string) model.Job {
	job := model.Job{}

	db.Connection.Where("status = ? AND type = ?", model.JobPending, jobType).First(&job)

	return job
}

// CountJobs count jobs by status
func (db *Database) CountJobs(status string) int {
	count := 0

	db.Connection.Model(&model.Job{}).Where("status = ?", status).Count(&count)

	return count
}

// DeleteJobByID deletes a job by id
func (db *Database) DeleteJobByID(id int) {
	db.Connection.Unscoped().Where("id=?", id).Delete(&migration.Job{})
}

// DeleteJobByUUID deletes a job by uuid
func (db *Database) DeleteJobByUUID(uuid string) {
	db.Connection.Unscoped().Where("uuid=?", uuid).Delete(&migration.Job{})
}

// UpdateJobByID updates a job by ID
func (db *Database) UpdateJobByID(job *model.Job) {
	db.Connection.Save(&job)
}

// CreateHost creates a new host
func (db *Database) CreateHost(host *model.Host) *model.Host {
	db.Connection.Create(host)
	return host
}

// HostExistByID check if host exists
func (db *Database) HostExistByID(id int) bool {
	host := model.Host{}

	db.Connection.Where("id = ?", id).First(&host)

	return host.ID > 0
}

// GetHostByID gets a host by id
func (db *Database) GetHostByID(id int) model.Host {
	host := model.Host{}

	db.Connection.Where("id = ?", id).First(&host)

	return host
}

// GetHosts gets services
func (db *Database) GetHosts() []model.Host {
	services := []model.Host{}

	db.Connection.Select("*").Find(&services)

	return services
}

// HostExistByUUID check if host exists
func (db *Database) HostExistByUUID(uuid string) bool {
	host := model.Host{}

	db.Connection.Where("uuid = ?", uuid).First(&host)

	return host.ID > 0
}

// GetHostByUUID gets a host by uuid
func (db *Database) GetHostByUUID(uuid string) model.Host {
	host := model.Host{}

	db.Connection.Where("uuid = ?", uuid).First(&host)

	return host
}

// DeleteHostByID deletes a host by id
func (db *Database) DeleteHostByID(id int) {
	db.Connection.Unscoped().Where("id=?", id).Delete(&migration.Host{})
	db.Connection.Unscoped().Where("host_id=?", id).Delete(&migration.Job{})
	db.Connection.Unscoped().Where("host_id=?", id).Delete(&migration.HostMeta{})
}

// DeleteHostByUUID deletes a host by uuid
func (db *Database) DeleteHostByUUID(uuid string) {
	host := db.GetHostByUUID(uuid)

	if host.ID > 0 {
		db.DeleteHostByID(host.ID)
	}
}

// UpdateHostByID updates a host by ID
func (db *Database) UpdateHostByID(host *model.Host) {
	db.Connection.Save(&host)
}

// CreateOption creates a new option
func (db *Database) CreateOption(option *model.Option) *model.Option {
	db.Connection.Create(option)
	return option
}

// OptionExistByID check if option exists
func (db *Database) OptionExistByID(id int) bool {
	option := model.Option{}

	db.Connection.Where("id = ?", id).First(&option)

	return option.ID > 0
}

// OptionExistByKey check if an option exists
func (db *Database) OptionExistByKey(key string) bool {
	option := model.Option{}

	db.Connection.Where("key = ?", key).First(&option)

	return option.ID > 0
}

// GetOptionByID gets an option by id
func (db *Database) GetOptionByID(id int) model.Option {
	option := model.Option{}

	db.Connection.Where("id = ?", id).First(&option)

	return option
}

// GetOptionByKey gets an option by key
func (db *Database) GetOptionByKey(key string) model.Option {
	option := model.Option{}

	db.Connection.Where("key = ?", key).First(&option)

	return option
}

// GetOptions gets options
func (db *Database) GetOptions() []model.Option {
	options := []model.Option{}

	db.Connection.Select("*").Find(&options)

	return options
}

// DeleteOptionByID deletes an option by id
func (db *Database) DeleteOptionByID(id int) {
	db.Connection.Unscoped().Where("id=?", id).Delete(&migration.Option{})
}

// DeleteOptionByKey deletes an option by key
func (db *Database) DeleteOptionByKey(key string) {
	db.Connection.Unscoped().Where("key=?", key).Delete(&migration.Option{})
}

// UpdateOptionByID updates an option by ID
func (db *Database) UpdateOptionByID(option *model.Option) {
	db.Connection.Save(&option)
}

// Close closes MySQL database connection
func (db *Database) Close() error {
	return db.Connection.Close()
}

// ReleaseChildJobs count jobs by status
func (db *Database) ReleaseChildJobs(parentID int) {
	db.Connection.Model(&model.Job{}).Where(
		"parent = ? AND status = ?",
		parentID,
		model.JobOnHold,
	).Update("status", model.JobPending)
}
