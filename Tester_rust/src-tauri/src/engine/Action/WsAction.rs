
use super::action::ActionObj;

use rusqlite::{params_from_iter, Connection};
use std::collections::HashMap;
use std::error::Error;

pub struct WSOpts {
    pub url: String,
    pub headers: HashMap<String, String>,
    pub query_params: HashMap<String, String>,
    pub subprotocols: Vec<String>,
    pub handshake_timeout: i64,
    pub write_timeout: i64,
}

pub struct WSMessage {
    pub message_id: String,
    pub sequence: i32,
    pub message_type: String,
    pub payload: Vec<u8>,
}

pub struct WSAction {
    pub action_id: String,
    pub app_id: String,
    pub action_name: String,
    pub environment_id: String,
    pub messages: Vec<WSMessage>,
    pub opts: WSOpts,
}

pub struct WSActionRepo {
    pub action_id: String,
    pub app_id: String,
    pub action_name: String,
    pub environment_id: String,
    pub opts: Vec<u8>,
}

impl ActionObj {

    // CREATE WS ACTION
    pub fn create_ws_action(&self, db: &Connection, app_id: &str, action_name: &str, _executor_type: &str, opts: &[u8], environment_id: &str) -> Result<WSActionRepo, Box<dyn Error>> {
        let action_id = uuid::Uuid::new_v4().to_string();

        let query = "
            INSERT INTO ws_actions
            (action_id, app_id, action_name, opts, environment_id)
            VALUES (?1, ?2, ?3, ?4, ?5)
        ";

        db.execute(query, (&action_id, app_id, action_name, opts, environment_id))?;

        Ok(WSActionRepo {
            action_id,
            app_id: app_id.to_string(),
            action_name: action_name.to_string(),
            environment_id: environment_id.to_string(),
            opts: opts.to_vec(),
        })
    }

    // GET WS ACTION
    pub fn get_ws_action(&self, db: &Connection, action_id: &str) -> Result<WSActionRepo, Box<dyn Error>> {
        let query = "
            SELECT action_id, app_id, action_name, opts, environment_id
            FROM ws_actions
            WHERE action_id = ?1
        ";

        let (action_id, app_id, action_name, opts, environment_id) =
            db.query_row(query, [action_id], |row| {
                Ok((
                    row.get::<_, String>(0)?,
                    row.get::<_, String>(1)?,
                    row.get::<_, String>(2)?,
                    row.get::<_, Vec<u8>>(3)?,
                    row.get::<_, String>(4)?,
                ))
            })?;

        Ok(WSActionRepo {
            action_id,
            app_id,
            action_name,
            opts,
            environment_id,
        })
    }

    // GET MULTIPLE WS ACTIONS
    pub fn get_ws_actions(&self, db: &Connection, ids: Vec<String>) -> Result<Vec<WSActionRepo>, Box<dyn Error>> {
        if ids.is_empty() {
            return Ok(Vec::new());
        }

        let placeholders = (1..=ids.len())
            .map(|i| format!("?{}", i))
            .collect::<Vec<String>>()
            .join(", ");

        let query = format!(
            "SELECT action_id, app_id, action_name, opts, environment_id
             FROM ws_actions
             WHERE action_id IN ({})",
            placeholders
        );

        let mut stmt = db.prepare(&query)?;

        let rows = stmt.query_map(params_from_iter(ids.iter()), |row| {
            Ok(WSActionRepo {
                action_id: row.get(0)?,
                app_id: row.get(1)?,
                action_name: row.get(2)?,
                opts: row.get(3)?,
                environment_id: row.get(4)?,
            })
        })?;

        let mut actions = Vec::new();

        for row in rows {
            actions.push(row?);
        }

        Ok(actions)
    }

    // UPDATE WS ACTION
    pub fn update_ws_action(&self, db: &Connection, action_id: &str, app_id: &str, action_name: &str, executor_type: &str, opts: &[u8], environment_id: &str) -> Result<(), Box<dyn Error>> {
        let query = "
            UPDATE ws_actions
            SET app_id = ?1,
                action_name = ?2,
                executor_type = ?3,
                opts = ?4,
                environment_id = ?5
            WHERE action_id = ?6
        ";

        let rows_affected = db.execute(
            query,
            (app_id, action_name, executor_type, opts, environment_id, action_id),
        )?;

        if rows_affected == 0 {
            return Err("WebSocket action not found".into());
        }

        Ok(())
    }

    // DELETE WS ACTION
    pub fn delete_ws_action(&self, db: &Connection, action_id: &str) -> Result<(), Box<dyn Error>> {
        let query = "DELETE FROM ws_actions WHERE action_id = ?1";

        db.execute(query, [action_id])?;

        Ok(())
    }
}
