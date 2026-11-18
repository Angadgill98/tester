
use std::collections::HashMap;
use std::error::Error;
use std::io::{Error as IoError, ErrorKind};
use rusqlite::{params_from_iter, Connection};

use crate::engine::Action;

pub struct HttpAction {
    pub action_id: String,
    pub action_name: String,
    pub app_id: String,
    pub environment_id: String,
    pub opts: HttpActionOpts,
}

pub struct HttpActionOpts {
    pub method: String,
    pub url: String,
    pub headers: HashMap<String, Vec<String>>,
    pub query: HashMap<String, Vec<String>>,
    pub body: Vec<u8>,
    pub cookies: HashMap<String, Vec<String>>,
}

pub struct HttpActionRepo {
    pub action_id: String,
    pub action_name: String,
    pub app_id: String,
    pub environment_id: String,
    pub opts: Vec<u8>,
}

impl Action::action::ActionObj {

    // CREATE HTTP ACTION
    pub fn create_http_action(&self, db: &Connection, action_name: &str, app_id: &str, opts: &[u8], environment_id: &str) -> Result<HttpActionRepo, Box<dyn Error>> {
        let action_id = uuid::Uuid::new_v4().to_string();

        let query = "
            INSERT INTO http_actions
            (action_id, action_name, app_id, opts, environment_id)
            VALUES (?1, ?2, ?3, ?4, ?5)
        ";

        db.execute(query, (&action_id, action_name, app_id, opts, environment_id))?;

        Ok(HttpActionRepo {
            action_id,
            action_name: action_name.to_string(),
            app_id: app_id.to_string(),
            environment_id: environment_id.to_string(),
            opts: opts.to_vec(),
        })
    }

    // GET HTTP ACTION
    pub fn get_http_action(&self, db: &Connection, action_id: &str) -> Result<HttpActionRepo, Box<dyn Error>> {
        let query = "
            SELECT action_id, app_id, action_name, opts, environment_id
            FROM http_actions
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

        Ok(HttpActionRepo {
            action_id,
            app_id,
            action_name,
            opts,
            environment_id,
        })
    }

    // GET MULTIPLE HTTP ACTIONS
    pub fn get_http_actions(&self, db: &Connection, ids: Vec<String>) -> Result<Vec<HttpActionRepo>, Box<dyn Error>> {
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
            Ok(HttpActionRepo {
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

    // UPDATE HTTP ACTION
    pub fn update_http_action(&self, db: &Connection, action_id: &str, action_name: &str, app_id: &str, opts: &[u8], environment_id: &str) -> Result<(), Box<dyn Error>> {
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
            (action_name, app_id, opts, environment_id, action_id),
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
    pub fn delete_http_action(&self, db: &Connection, action_id: &str) -> Result<(), Box<dyn Error>> {
        let query = "DELETE FROM http_actions WHERE action_id = ?1";

        db.execute(query, [action_id])?;

        Ok(())
    }
}