package postgres

import (
	"context"
	"database/sql"

	"simple/constants"
)

// MigrateSchema creates the required database schema for PostgreSQL
// This includes tables and indexes needed for the application
func MigrateSchema(db *sql.DB) error {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DBQueryTimeout)
	defer cancel()

	// Create records table if not exists
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS records (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return err
	}

	// Create index on created_at for efficient sorting
	_, err = db.ExecContext(ctx, `
		CREATE INDEX IF NOT EXISTS idx_records_created_at ON records(created_at DESC)
	`)
	if err != nil {
		return err
	}

	return nil
}

// GetSchemaVersion returns the current schema version
func GetSchemaVersion() string {
	return "v1"
}