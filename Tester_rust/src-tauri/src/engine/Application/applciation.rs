
use rusqlite::{params_from_iter, Connection};
use serde::{Deserialize, Serialize};
use std::error::Error;

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Application {
    pub app_id: String,
    pub app_name: String,
    pub parent_app_id: String,
    pub environment_id: Option<String>,
    pub actions: Vec<String>,
    pub sequences: Vec<String>,
    pub children: Vec<String>,
}

pub struct ApplicationObj;

impl ApplicationObj {
    pub fn create_application_obj() -> Self {
        Self {}
    }

    pub fn create_application(&self, db: &Connection, app_name: &str, parent_app_id: String) -> Result<Application, Box<dyn Error>> {
        let app_id = uuid::Uuid::new_v4().to_string();
    
        let actions: Vec<String> = Vec::new();
        let sequences: Vec<String> = Vec::new();
        let children: Vec<String> = Vec::new();
    
        let actions_json = serde_json::to_string(&actions)?;
        let sequences_json = serde_json::to_string(&sequences)?;
        let children_json = serde_json::to_string(&children)?;
    
        let query = "
            INSERT INTO applications
            (id, name, parent_app_id, action_ids, sequence_ids, children)
            VALUES (?1, ?2, ?3, ?4, ?5, ?6)
        ";
    
        db.execute(query, (&app_id, app_name, &parent_app_id, actions_json, sequences_json, children_json))?;
    
        Ok(Application {
            app_id,
            app_name: app_name.to_string(),
            parent_app_id,
            environment_id: None,
            actions,
            sequences,
            children,
        })
    }
    
    pub fn get_application_by_id(&self, db: &Connection, app_id: String) -> Result<Application, Box<dyn Error>> {
        let query = "
            SELECT id, name, parent_app_id, environment_id,
                   action_ids, sequence_ids, children
            FROM applications
            WHERE id = ?1
        ";

        let (app_id, app_name, parent_app_id, environment_id, actions_json, sequences_json, children_json) =
            db.query_row(query, [&app_id], |row| {
                Ok((
                    row.get::<_, String>(0)?,
                    row.get::<_, String>(1)?,
                    row.get::<_, String>(2)?,
                    row.get::<_, Option<String>>(3)?,
                    row.get::<_, String>(4)?,
                    row.get::<_, String>(5)?,
                    row.get::<_, String>(6)?,
                ))
            })?;

        let actions = serde_json::from_str(&actions_json)?;
        let sequences = serde_json::from_str(&sequences_json)?;
        let children = serde_json::from_str(&children_json)?;

        Ok(Application {
            app_id,
            app_name,
            parent_app_id,
            environment_id,
            actions,
            sequences,
            children,
        })
    }

    pub fn get_applications_by_ids(&self, db: &Connection, app_ids: Vec<String>) -> Result<Vec<Application>, Box<dyn Error>> {
        if app_ids.is_empty() {
            return Ok(Vec::new());
        }

        let placeholders = (1..=app_ids.len())
            .map(|i| format!("?{}", i))
            .collect::<Vec<String>>()
            .join(", ");

        let query = format!(
            "SELECT id, name, parent_app_id, environment_id,
                    action_ids, sequence_ids, children
             FROM applications
             WHERE id IN ({})",
            placeholders
        );

        let mut stmt = db.prepare(&query)?;

        let rows = stmt.query_map(params_from_iter(app_ids.iter()), |row| {
            Ok((
                row.get::<_, String>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, String>(2)?,
                row.get::<_, Option<String>>(3)?,
                row.get::<_, String>(4)?,
                row.get::<_, String>(5)?,
                row.get::<_, String>(6)?,
            ))
        })?;

        let mut applications = Vec::new();

        for row in rows {
            let (app_id, app_name, parent_app_id, environment_id, actions_json, sequences_json, children_json) = row?;

            applications.push(Application {
                app_id,
                app_name,
                parent_app_id,
                environment_id,
                actions: serde_json::from_str(&actions_json)?,
                sequences: serde_json::from_str(&sequences_json)?,
                children: serde_json::from_str(&children_json)?,
            });
        }

        Ok(applications)
    }

    pub fn register_action(&self, db: &Connection, app_id: String, action_id: String) -> Result<(), Box<dyn Error>> {
        db.execute(
            "UPDATE applications SET action_ids = json_insert(action_ids, '$[#]', ?1) WHERE id = ?2",
            (action_id, app_id),
        )?;

        Ok(())
    }

    pub fn register_sequence(&self, db: &Connection, app_id: String, sequence_id: String) -> Result<(), Box<dyn Error>> {
        db.execute(
            "UPDATE applications SET sequence_ids = json_insert(sequence_ids, '$[#]', ?1) WHERE id = ?2",
            (sequence_id, app_id),
        )?;

        Ok(())
    }

    pub fn register_child(&self, db: &Connection, app_id: String, child_id: String) -> Result<(), Box<dyn Error>> {
        db.execute(
            "UPDATE applications SET children = json_insert(children, '$[#]', ?1) WHERE id = ?2",
            (child_id, app_id),
        )?;

        Ok(())
    }
}
