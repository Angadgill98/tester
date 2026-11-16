use super::action::ActionObj;

use rusqlite::{params_from_iter, Connection};
use serde::{Deserialize, Serialize};
use serde_json::value::RawValue;
use std::collections::HashMap;
use std::error::Error;

#[derive(Debug, Serialize, Deserialize)]
pub struct WSOpts {
    pub url: String,
    pub headers: HashMap<String, String>,
    pub query_params: HashMap<String, String>,
    pub subprotocols: Vec<String>,
    pub messages: Vec<WSMessage>,
    pub handshake_timeout: i64,
    pub write_timeout: i64,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct WSMessage {
    pub message_id: String,
    pub sequence: i32,
    pub message_type: String,
    pub payload: Box<RawValue>,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct WSAction {
    pub action_id: String,
    pub app_id: String,
    pub action_name: String,
    pub environment_id: String,
    pub opts: WSOpts,
}

impl ActionObj {
    pub fn create_ws_action(&self,db: &Connection, app_id: &str, action_name: &str, _executor_type: &str, opts: &WSOpts, environment_id: &str) -> Result<String, Box<dyn Error>> {
        let action_id = uuid::Uuid::new_v4().to_string();
        let opts_json = serde_json::to_string(opts)?;
    
        let query = "
            INSERT INTO ws_actions
            (action_id, app_id, action_name, opts, environment_id)
            VALUES (?1, ?2, ?3, ?4, ?5)
        ";
    
        db.execute(query, (&action_id, app_id, action_name, opts_json, environment_id))?;
    
        Ok(action_id)
    }
    
    pub fn get_ws_action(&self,db: &Connection, action_id: &str) -> Result<WSAction, Box<dyn Error>> {
        let query = "
            SELECT action_id, app_id, action_name, opts, environment_id
            FROM ws_actions
            WHERE action_id = ?1
        ";

        let (action_id, app_id, action_name, opts_json, environment_id) =
            db.query_row(query, [action_id], |row| {
                Ok((
                    row.get::<_, String>(0)?,
                    row.get::<_, String>(1)?,
                    row.get::<_, String>(2)?,
                    row.get::<_, String>(3)?,
                    row.get::<_, String>(4)?,
                ))
            })?;

        let opts: WSOpts = serde_json::from_str(&opts_json)?;

        Ok(WSAction {
            action_id,
            app_id,
            action_name,
            opts,
            environment_id,
        })
    }

    pub fn get_ws_actions(&self,db: &Connection, ids: Vec<String>) -> Result<Vec<WSAction>, Box<dyn Error>> {
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
            Ok((
                row.get::<_, String>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, String>(2)?,
                row.get::<_, String>(3)?,
                row.get::<_, String>(4)?,
            ))
        })?;

        let mut actions = Vec::new();

        for row in rows {
            let (action_id, app_id, action_name, opts_json, environment_id) = row?;

            let opts: WSOpts = serde_json::from_str(&opts_json)?;

            actions.push(WSAction {
                action_id,
                app_id,
                action_name,
                opts,
                environment_id,
            });
        }

        Ok(actions)
    }

    pub fn update_ws_action(
        &self,
        db: &Connection,
        action_id: &str,
        app_id: &str,
        action_name: &str,
        executor_type: &str,
        opts: &str,
        environment_id: &str,
    ) -> Result<(), Box<dyn Error>> {
        // Validate the JSON before updating.
        let _: serde_json::Value = serde_json::from_str(opts)?;

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

    pub fn delete_ws_action(&self,db: &Connection, action_id: &str) -> Result<(), Box<dyn Error>> {
        let query = "DELETE FROM ws_actions WHERE action_id = ?1";

        db.execute(query, [action_id])?;

        Ok(())
    }
}