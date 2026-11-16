use std::collections::HashMap;

use serde::{Deserialize, Serialize};
use serde_json::value::RawValue;

use rusqlite::{Connection, params_from_iter};
use std::error::Error;
use std::io::{Error as IoError, ErrorKind};

use crate::engine::Action;

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct HttpAction {
    pub action_id: String,
    pub action_name: String,
    pub app_id: String,
    pub environment_id: String,
    pub opts: HttpActionOpts,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct HttpActionOpts {
    pub method: String,
    pub url: String,
    pub headers: HashMap<String, Vec<String>>,
    pub query: HashMap<String, Vec<String>>,
    pub body: Box<RawValue>,
    pub cookies: HashMap<String, Vec<String>>,
}

impl Action::action::ActionObj {

    // CREATE HTTP ACTION
    pub fn create_http_action(&self,db: &Connection, action_name: &str, app_id: &str, opts: &HttpActionOpts, environment_id: &str) -> Result<String, Box<dyn Error>> {
        let action_id = uuid::Uuid::new_v4().to_string();
        let opts_json = serde_json::to_string(opts)?;
    
        let query = "
            INSERT INTO http_actions
            (action_id, action_name, app_id, opts, environment_id)
            VALUES (?1, ?2, ?3, ?4, ?5)
        ";
    
        db.execute(query, (&action_id, action_name, app_id, opts_json, environment_id))?;
    
        Ok(action_id)
    }

    // GET HTTP ACTION
    pub fn get_http_action(&self,db: &Connection, action_id: &str) -> Result<HttpAction, Box<dyn Error>> {
        let query = "
            SELECT action_id, app_id, action_name, opts, environment_id
            FROM http_actions
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

        let opts: HttpActionOpts = serde_json::from_str(&opts_json)?;

        Ok(HttpAction {
            action_id,
            app_id,
            action_name,
            opts,
            environment_id,
        })
    }

    // GET MULTIPLE HTTP ACTIONS
    pub fn get_http_actions(&self,db: &Connection, ids: Vec<String>) -> Result<Vec<HttpAction>, Box<dyn Error>> {
        if ids.is_empty() {
            return Ok(Vec::new());
        }

        let placeholders = (1..=ids.len())
            .map(|i| format!("?{}", i))
            .collect::<Vec<String>>()
            .join(", ");

        let query = format!(
            "SELECT action_id, app_id, action_name, opts, environment_id
             FROM http_actions
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

            let opts: HttpActionOpts = serde_json::from_str(&opts_json)?;

            actions.push(HttpAction {
                action_id,
                app_id,
                action_name,
                opts,
                environment_id,
            });
        }

        Ok(actions)
    }

    // UPDATE HTTP ACTION
    pub fn update_http_action(
        &self,
        db: &Connection,
        action_id: &str,
        action_name: &str,
        app_id: &str,
        opts: &HttpActionOpts,
        environment_id: &str,
    ) -> Result<(), Box<dyn Error>> {
        let opts_json = serde_json::to_string(opts)?;

        let query = "
            UPDATE http_actions
            SET action_name = ?1,
                app_id = ?2,
                opts = ?3,
                environment_id = ?4
            WHERE action_id = ?5
        ";

        let rows_affected = db.execute(
            query,
            (action_name, app_id, opts_json, environment_id, action_id),
        )?;

        if rows_affected == 0 {
            return Err(Box::new(IoError::new(
                ErrorKind::NotFound,
                "HTTP action not found",
            )));
        }

        Ok(())
    }

    // DELETE HTTP ACTION
    pub fn delete_http_action(&self,db: &Connection, action_id: &str) -> Result<(), Box<dyn Error>> {
        let query = "DELETE FROM http_actions WHERE action_id = ?1";

        db.execute(query, [action_id])?;

        Ok(())
    }
}