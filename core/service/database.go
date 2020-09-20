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
	db.Connection.AutoMigrate(&migration.Job{}, &migration.Agent{}, &migration.Host{})
	status = status && db.Connection.HasTable(&migration.Job{})
	status = status && db.Connection.HasTable(&migration.Agent{})
	status = status && db.Connection.HasTable(&migration.Host{})
	return status
}

// Rollback drop tables
func (db *Database) Rollback() bool {
	status := true
	db.Connection.DropTableIfExists(&migration.Job{}, &migration.Agent{}, &migration.Host{})
	status = status && !db.Connection.HasTable(&migration.Job{})
	status = status && !db.Connection.HasTable(&migration.Agent{})
	status = status && !db.Connection.HasTable(&migration.Host{})
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

// CreateHost creates a new service
func (db *Database) CreateHost(service *model.Host) *model.Host {
	db.Connection.Create(service)
	return service
}

// HostExistByID check if service exists
func (db *Database) HostExistByID(id int) bool {
	service := model.Host{}

	db.Connection.Where("id = ?", id).First(&service)

	return service.ID > 0
}

// GetHostByID gets a service by id
func (db *Database) GetHostByID(id int) model.Host {
	service := model.Host{}

	db.Connection.Where("id = ?", id).First(&service)

	return service
}

// GetHosts gets services
func (db *Database) GetHosts() []model.Host {
	services := []model.Host{}

	db.Connection.Select("*").Find(&services)

	return services
}

// HostExistByUUID check if service exists
func (db *Database) HostExistByUUID(uuid string) bool {
	service := model.Host{}

	db.Connection.Where("uuid = ?", uuid).First(&service)

	return service.ID > 0
}

// GetHostByUUID gets a service by uuid
func (db *Database) GetHostByUUID(uuid string) model.Host {
	service := model.Host{}

	db.Connection.Where("uuid = ?", uuid).First(&service)

	return service
}

// DeleteHostByID deletes a service by id
func (db *Database) DeleteHostByID(id int) {
	db.Connection.Unscoped().Where("id=?", id).Delete(&migration.Host{})
}

// DeleteHostByUUID deletes a service by uuid
func (db *Database) DeleteHostByUUID(uuid string) {
	db.Connection.Unscoped().Where("uuid=?", uuid).Delete(&migration.Host{})
}

// UpdateHostByID updates a service by ID
func (db *Database) UpdateHostByID(service *model.Host) {
	db.Connection.Save(&service)
}

// CreateAgent creates a new agent
func (db *Database) CreateAgent(agent *model.Agent) *model.Agent {
	db.Connection.Create(agent)
	return agent
}

// AgentExistByID check if agent exists
func (db *Database) AgentExistByID(id int) bool {
	agent := model.Agent{}

	db.Connection.Where("id = ?", id).First(&agent)

	return agent.ID > 0
}

// GetAgentByID gets a agent by id
func (db *Database) GetAgentByID(id int) model.Agent {
	agent := model.Agent{}

	db.Connection.Where("id = ?", id).First(&agent)

	return agent
}

// GetAgents gets agents
func (db *Database) GetAgents() []model.Agent {
	agents := []model.Agent{}

	db.Connection.Select("*").Find(&agents)

	return agents
}

// AgentExistByUUID check if agent exists
func (db *Database) AgentExistByUUID(uuid string) bool {
	agent := model.Agent{}

	db.Connection.Where("uuid = ?", uuid).First(&agent)

	return agent.ID > 0
}

// GetAgentByUUID gets a agent by uuid
func (db *Database) GetAgentByUUID(uuid string) model.Agent {
	agent := model.Agent{}

	db.Connection.Where("uuid = ?", uuid).First(&agent)

	return agent
}

// DeleteAgentByID deletes a agent by id
func (db *Database) DeleteAgentByID(id int) {
	db.Connection.Unscoped().Where("id=?", id).Delete(&migration.Agent{})
}

// DeleteAgentByUUID deletes a agent by uuid
func (db *Database) DeleteAgentByUUID(uuid string) {
	db.Connection.Unscoped().Where("uuid=?", uuid).Delete(&migration.Agent{})
}

// UpdateAgentByID updates a agent by ID
func (db *Database) UpdateAgentByID(agent *model.Agent) {
	db.Connection.Save(&agent)
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
