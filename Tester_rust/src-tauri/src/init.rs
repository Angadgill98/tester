
use std::path::Path;

use rusqlite::{Connection, Result};

pub fn init_db(db_path: impl AsRef<Path>) -> Result<Connection> {
    let db = Connection::open(db_path)?;

    db.execute_batch("PRAGMA foreign_keys = ON;")?;

    init_tables(&db)?;
    initialize_global_env(&db)?;
    initialize_global_app(&db)?;

    Ok(db)
}

pub fn initialize_global_app(db: &Connection) -> Result<()> {
    let query = "
        INSERT OR IGNORE INTO applications
        (id, name, parent_app_id, environment_id, action_ids, sequence_ids, children)
        VALUES ('0', 'GLOBAL', '0', '0', '[]', '[]', '[]')
    ";

    db.execute(query, [])?;
    Ok(())
}

pub fn initialize_global_env(db: &Connection) -> Result<()> {
    let query = "
        INSERT OR IGNORE INTO environments
        (environment_id, environment_name, keys_values)
        VALUES ('0', 'GLOBAL', '{}')
    ";

    db.execute(query, [])?;
    Ok(())
}

pub fn setup_applications_table(db: &Connection) -> Result<()> {
    let query = "
        CREATE TABLE IF NOT EXISTS applications (
            id TEXT PRIMARY KEY,
            name TEXT NOT NULL,
            parent_app_id TEXT,
            environment_id TEXT,
            action_ids TEXT NOT NULL DEFAULT '[]',
            sequence_ids TEXT NOT NULL DEFAULT '[]',
            children TEXT NOT NULL DEFAULT '[]'
        )
    ";

    db.execute(query, [])?;
    Ok(())
}

pub fn create_actions_table(db: &Connection) -> Result<()> {
    let query = "
        CREATE TABLE IF NOT EXISTS actions (
            action_id TEXT PRIMARY KEY,
            action_type TEXT NOT NULL
        )
    ";

    db.execute(query, [])?;
    Ok(())
}

pub fn create_http_actions_table(db: &Connection) -> Result<()> {
    let query = "
        CREATE TABLE IF NOT EXISTS http_actions (
            action_id TEXT PRIMARY KEY,
            action_name TEXT NOT NULL,
            app_id TEXT NOT NULL,
            environment_id TEXT,
            opts TEXT NOT NULL DEFAULT '{}'
        )
    ";

    db.execute(query, [])?;
    Ok(())
}

pub fn create_ws_actions_table(db: &Connection) -> Result<()> {
    let query = "
        CREATE TABLE IF NOT EXISTS ws_actions (
            action_id TEXT PRIMARY KEY,
            action_name TEXT NOT NULL,
            app_id TEXT NOT NULL,
            environment_id TEXT,
            opts TEXT NOT NULL DEFAULT '{}'
        )
    ";

    db.execute(query, [])?;
    Ok(())
}

pub fn create_environments_table(db: &Connection) -> Result<()> {
    let query = "
        CREATE TABLE IF NOT EXISTS environments (
            environment_id TEXT PRIMARY KEY,
            environment_name TEXT NOT NULL,
            keys_values TEXT NOT NULL DEFAULT '{}'
        )
    ";

    db.execute(query, [])?;
    Ok(())
}

pub fn create_sequences_table(db: &Connection) -> Result<()> {
    let query = "
        CREATE TABLE IF NOT EXISTS sequences (
            sequence_id TEXT PRIMARY KEY,
            app_id TEXT NOT NULL,
            sequence_name TEXT NOT NULL,
            action_ids TEXT NOT NULL DEFAULT '[]'
        )
    ";

    db.execute(query, [])?;
    Ok(())
}

pub fn init_tables(db: &Connection) -> Result<()> {
    create_environments_table(db)?;
    setup_applications_table(db)?;
    create_actions_table(db)?;
    create_http_actions_table(db)?;
    create_ws_actions_table(db)?;
    create_sequences_table(db)?;

    Ok(())
}
