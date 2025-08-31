// Copyright 2025 Baleine Jay
// Licensed under the Phicode Non-Commercial License (https://banes-lab.com/licensing)
// Commercial use requires a paid license. See link for details.
package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/glebarez/go-sqlite"
)

const schema = `
CREATE TABLE IF NOT EXISTS benchmarks (
	id INTEGER PRIMARY KEY,
	repo TEXT NOT NULL,
	branch TEXT NOT NULL,
	commit_hash TEXT NOT NULL,
	component TEXT NOT NULL,
	value REAL NOT NULL,
	timestamp INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS baselines (
	repo TEXT NOT NULL,
	component TEXT NOT NULL,
	baseline_value REAL NOT NULL,
	sample_count INTEGER DEFAULT 5,
	updated_at INTEGER NOT NULL,
	PRIMARY KEY (repo, component)
);

CREATE TABLE IF NOT EXISTS config (
	repo TEXT PRIMARY KEY,
	threshold_percent REAL DEFAULT 10.0,
	min_samples INTEGER DEFAULT 5,
	enabled BOOLEAN DEFAULT 1
);

CREATE INDEX IF NOT EXISTS idx_benchmarks_repo_component ON benchmarks(repo, component);
CREATE INDEX IF NOT EXISTS idx_benchmarks_timestamp ON benchmarks(timestamp);
`

func Init(path string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err := createTables(db); err != nil {
		return nil, fmt.Errorf("table creation failed: %w", err)
	}

	return db, nil
}

func createTables(db *sqlx.DB) error {
	_, err := db.Exec(schema)
	return err
}