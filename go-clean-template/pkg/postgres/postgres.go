// Package postgres implements postgres connection.
package postgres

import (
	"database/sql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// Postgres -.
type Postgres struct {
	*gorm.DB
}

// New -.
func New(url string) (*gorm.DB, *sql.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
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
