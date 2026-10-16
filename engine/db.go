package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

func SetSQLite() (*sql.DB, error) {
	if err := os.MkdirAll("data", 0755); err != nil {
		return nil, fmt.Errorf("create data directory: %w", err)
	}

	db, err := sql.Open("sqlite", "data/app.db")
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	if err := SetupAppsTable(db); err != nil {
		db.Close()
		return nil, err
	}

	if err := SetupHTTPActionsTable(db); err != nil {
		db.Close()
		return nil, err
	}

	if err := SetupWSActionsTable(db); err != nil {
		db.Close()
		return nil, err
	}

	if err := SetupSequencesTable(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func SetupAppsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS apps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			action_ids TEXT NOT NULL DEFAULT '[]',
			sequence_ids TEXT NOT NULL DEFAULT '[]'
		)
	`)
	if err != nil {
		return fmt.Errorf("create apps table: %w", err)
	}

	return nil
}

func SetupHTTPActionsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS http_actions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			app_id INTEGER NOT NULL,
			action_name TEXT NOT NULL,
			opts BLOB NOT NULL DEFAULT X'',
			FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
			UNIQUE (app_id, action_name)
		)
	`)
	if err != nil {
		return fmt.Errorf("create http_actions table: %w", err)
	}

	return nil
}

func SetupWSActionsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS ws_actions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			app_id INTEGER NOT NULL,
			action_name TEXT NOT NULL,
			executor_type TEXT NOT NULL,
			opts BLOB NOT NULL DEFAULT X'',
			FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
			UNIQUE (app_id, action_name)
		)
	`)
	if err != nil {
		return fmt.Errorf("create ws_actions table: %w", err)
	}

	return nil
}

func SetupSequencesTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS sequences (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			app_id INTEGER NOT NULL,
			sequence_name TEXT NOT NULL,
			action_ids TEXT NOT NULL DEFAULT '[]',
			FOREIGN KEY (app_id) REFERENCES apps(id) ON DELETE CASCADE,
			UNIQUE (app_id, sequence_name)
		)
	`)
	if err != nil {
		return fmt.Errorf("create sequences table: %w", err)
	}

	return nil
}