//go:build sqlite

package sqlite

import (
	"database/sql"
	"fmt"
)

const schemaVersion = 1

// migrate creates or upgrades the database schema.
func (s *SQLiteProvider) migrate() error {
	// Create schema_version table if it doesn't exist
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_version (
			version INTEGER NOT NULL
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create schema_version table: %w", err)
	}

	// Check current version
	var currentVersion int
	err = s.db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_version").Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to read schema version: %w", err)
	}

	if currentVersion < schemaVersion {
		if err := s.applyMigrations(currentVersion); err != nil {
			return err
		}
	}

	return nil
}

func (s *SQLiteProvider) applyMigrations(fromVersion int) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin migration transaction: %w", err)
	}
	defer tx.Rollback()

	if fromVersion < 1 {
		if err := createV1Schema(tx); err != nil {
			return err
		}
		if _, err := tx.Exec("INSERT INTO schema_version (version) VALUES (1)"); err != nil {
			return fmt.Errorf("failed to record schema version 1: %w", err)
		}
	}

	return tx.Commit()
}

func createV1Schema(tx *sql.Tx) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS prompts (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL UNIQUE,
			title TEXT,
			description TEXT,
			arguments TEXT,
			content TEXT NOT NULL,
			tags TEXT,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS idx_prompts_name ON prompts(name)`,
		`CREATE INDEX IF NOT EXISTS idx_prompts_tags ON prompts(tags)`,
	}

	for _, stmt := range statements {
		if _, err := tx.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	return nil
}
