package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)







func InitDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if err := InitTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func InitializeGlobalApp(db *sql.DB) error {
	query := `
		INSERT OR IGNORE INTO applications (id, name, parent_app_id, action_ids, sequence_ids, children)
		VALUES (0, 'GLOBAL', 0, '[]', '[]', '[]')
	`

	_, err := db.Exec(query)
	return err
}
func InitializeGlobalEnv(db *sql.DB) error {
    query := `INSERT OR IGNORE INTO environments (environment_id, environment_name, keys_values) VALUES (0, 'GLOBAL', '{}')`

    _, err := db.Exec(query)
    if err != nil {
        return err
    }

    query = `UPDATE applications SET environment_id = 0 WHERE id = 0`

    _, err = db.Exec(query)
    return err
}
func SetupApplicationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS applications (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			parent_app_id INTEGER,
			environment_id INTEGER,
			action_ids TEXT NOT NULL DEFAULT '[]',
			sequence_ids TEXT NOT NULL DEFAULT '[]',
			children TEXT NOT NULL DEFAULT '[]',
			FOREIGN KEY (parent_app_id) REFERENCES applications(id) ON DELETE CASCADE,
			FOREIGN KEY (environment_id) REFERENCES environments(environment_id)
		)
	`)
	if err != nil {
		return fmt.Errorf("create applications table: %w", err)
	}

	return nil
}

func CreateActionsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS actions (
		action_id INTEGER PRIMARY KEY AUTOINCREMENT,
		action_type TEXT NOT NULL 
	);`

	_, err := db.Exec(query)
	if err != nil { 
		return fmt.Errorf("create actions table: %w", err)
	}

	return nil
}

func CreateHTTPActionsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS http_actions (
		action_id INTEGER PRIMARY KEY AUTOINCREMENT,
		action_name TEXT NOT NULL,
		app_id INTEGER NOT NULL,
		environment_id INTEGER,
		opts TEXT NOT NULL DEFAULT '{}',
		FOREIGN KEY (action_id) REFERENCES actions(action_id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create http_actions table: %w", err)
	}

	return nil
}

func CreateWSActionsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS ws_actions (
		action_id INTEGER PRIMARY KEY AUTOINCREMENT,
		action_name TEXT NOT NULL,
		app_id INTEGER NOT NULL,
		environment_id INTEGER,
		opts TEXT NOT NULL DEFAULT '{}',
		FOREIGN KEY (action_id) REFERENCES actions(action_id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create ws_actions table: %w", err)
	}

	return nil
}

func CreateEnvironmentsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS environments (
		environment_id INTEGER PRIMARY KEY AUTOINCREMENT,
		environment_name TEXT NOT NULL,
		keys_values TEXT NOT NULL DEFAULT '{}'
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create environments table: %w", err)
	}

	return nil
}

func CreateSequencesTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS sequences (
		sequence_id INTEGER PRIMARY KEY AUTOINCREMENT,
		app_id INTEGER NOT NULL,
		sequence_name TEXT NOT NULL,
		action_ids TEXT NOT NULL,
		FOREIGN KEY (app_id) REFERENCES applications(id)
	);`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("create sequences table: %w", err)
	}

	return nil
}

func InitTables(db *sql.DB) error {
	if err := SetupApplicationsTable(db); err != nil {
		return err
	}

	if err := CreateActionsTable(db); err != nil {
		return err
	}

	if err := CreateHTTPActionsTable(db); err != nil {
		return err
	}

	if err := CreateWSActionsTable(db); err != nil {
		return err
	}

	if err := CreateEnvironmentsTable(db); err != nil {
		return err
	}

	if err := CreateSequencesTable(db); err != nil {
		return err
	}

	return nil
}