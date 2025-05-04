package repository

import (
	"database/sql"
	"fmt"
	"time"

	"library-management-system/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

// Database represents a database connection
type Database struct {
	*sql.DB
}

// DatabaseConfig represents database configuration
type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	Params   string
}

// NewDatabase creates a new Database instance
func NewDatabase(cfg config.DatabaseConfig) (*Database, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.Params)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening database connection: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetimeMinutes) * time.Minute)

	// Verify connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &Database{db}, nil
}

// Close closes the database connection
func (db *Database) Close() error {
	return db.DB.Close()
}

// Transaction executes a function within a database transaction
func (db *Database) Transaction(fn func(*sql.Tx) error) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // re-throw panic after rollback
		}
	}()

	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}
