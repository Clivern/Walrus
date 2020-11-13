// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"time"

	"github.com/clivern/walrus/core/migration"
)

var (
	// JobPending pending job type
	JobPending = "pending"

	// JobFailed failed job type
	JobFailed = "failed"

	// JobSuccess success job type
	JobSuccess = "success"

	// JobOnHold on hold job type
	JobOnHold = "on_hold"
)

// Job struct
type Job struct {
	ID        int        `json:"id"`
	UUID      string     `json:"uuid"`
	Payload   string     `json:"payload"`
	Status    string     `json:"status"`
	Type      string     `json:"type"`
	Result    string     `json:"result"`
	Retry     int        `json:"retry"`
	Parent    int        `json:"parent"`
	RunAt     *time.Time `json:"run_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Jobs struct
type Jobs struct {
	Jobs []Job `json:"jobs"`
}

// CreateJob creates a new job
func (db *Database) CreateJob(job *Job) *Job {
	db.Connection.Create(job)
	return job
}

// JobExistByID check if job exists
func (db *Database) JobExistByID(id int) bool {
	job := Job{}

	db.Connection.Where("id = ?", id).First(&job)

	return job.ID > 0
}

// GetJobByID gets a job by id
func (db *Database) GetJobByID(id int) Job {
	job := Job{}

	db.Connection.Where("id = ?", id).First(&job)

	return job
}

// GetJobs gets jobs
func (db *Database) GetJobs() []Job {
	jobs := []Job{}

	db.Connection.Select("*").Find(&jobs)

	return jobs
}

// JobExistByUUID check if job exists
func (db *Database) JobExistByUUID(uuid string) bool {
	job := Job{}

	db.Connection.Where("uuid = ?", uuid).First(&job)

	return job.ID > 0
}

// GetJobByUUID gets a job by uuid
func (db *Database) GetJobByUUID(uuid string) Job {
	job := Job{}

	db.Connection.Where("uuid = ?", uuid).First(&job)

	return job
}

// GetPendingJobByType gets a job by uuid
func (db *Database) GetPendingJobByType(jobType string) Job {
	job := Job{}

	db.Connection.Where("status = ? AND type = ?", JobPending, jobType).First(&job)

	return job
}

// CountJobs count jobs by status
func (db *Database) CountJobs(status string) int {
	count := 0

	db.Connection.Model(&Job{}).Where("status = ?", status).Count(&count)

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
func (db *Database) UpdateJobByID(job *Job) {
	db.Connection.Save(&job)
}
