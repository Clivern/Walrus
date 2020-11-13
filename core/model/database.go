// Copyright 2020 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package model

import (
	"fmt"
	"time"

	"github.com/clivern/walrus/core/migration"

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
func (db *Database) Connect(dsn DSN) error {
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

	dsn := DSN{
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

// Close closes MySQL database connection
func (db *Database) Close() error {
	return db.Connection.Close()
}

// ReleaseChildJobs count jobs by status
func (db *Database) ReleaseChildJobs(parentID int) {
	db.Connection.Model(&Job{}).Where(
		"parent = ? AND status = ?",
		parentID,
		JobOnHold,
	).Update("status", JobPending)
}
