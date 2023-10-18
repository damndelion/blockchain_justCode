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
func New(url string) (*sql.DB, *Postgres, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	sqlDB, err := db.DB()
	if err != nil {
		return nil, nil, err
	}

	return sqlDB, &Postgres{db}, nil
}

// Close -.
func (p *Postgres) Close() {
	d, err := p.DB.DB()
	if err != nil {
		log.Printf("error closing database: %s", err.Error())
	}

	d.Close()
}
