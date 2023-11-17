// Package postgres implements postgres connection.
package postgres

import (
	"database/sql"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Postgres -.
type Postgres struct {
	*gorm.DB
}

// New -.
func New(url string) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return db, sqlDB, nil
}

// Close -.
func (p *Postgres) Close() {
	sqlDB, err := p.DB.DB()
	if err != nil {
		log.Printf("error closing database: %s", err.Error())
	}

	sqlDB.Close()
}
