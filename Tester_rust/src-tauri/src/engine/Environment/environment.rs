
use rusqlite::Connection;
use serde::{Deserialize, Serialize};
use serde_json::value::RawValue;
use serde_json::Value;
use std::collections::HashMap;
use std::error::Error;

#[derive(Debug, Serialize, Deserialize)]
pub struct Environment {
    pub environment_id: String,
    pub environment_name: String,
    pub keys_values: HashMap<String, String>,
}

impl Environment {
    pub fn create_environment(db: &Connection, environment_name: &str) -> Result<Self, Box<dyn Error>> {
        let environment_id = uuid::Uuid::new_v4().to_string();
        let keys_values: HashMap<String, String> = HashMap::new();
    
        let keys_values_json = serde_json::to_string(&keys_values)?;
    
        let query = "
            INSERT INTO environments (environment_id, environment_name, keys_values)
            VALUES (?1, ?2, ?3)
        ";
    
        db.execute(query, (&environment_id, environment_name, keys_values_json))?;
    
        Ok(Self {
            environment_id,
            environment_name: environment_name.to_string(),
            keys_values,
        })
    }

    pub fn get_environment(db: &Connection, environment_id: String) -> Result<Self, Box<dyn Error>> {
        let query = "
            SELECT environment_id, environment_name, keys_values
            FROM environments
            WHERE environment_id = ?1
        ";
    
        let (environment_id, environment_name, keys_values_json) =
            db.query_row(query, [&environment_id], |row| {
                Ok((
                    row.get::<_, String>(0)?,
                    row.get::<_, String>(1)?,
                    row.get::<_, String>(2)?,
                ))
            })?;
    
        let keys_values = serde_json::from_str(&keys_values_json)?;
    
        Ok(Self {
            environment_id,
            environment_name,
            keys_values,
        })
    }
    
    pub fn update_environment(
        db: &Connection,
        environment_id: String,
        environment_name: &str,
        keys_values: &RawValue,
    ) -> Result<(), Box<dyn Error>> {
        let query = "
            UPDATE environments
            SET environment_name = ?1,
                keys_values = ?2
            WHERE environment_id = ?3
        ";

        let rows_affected = db.execute(
            query,
            (environment_name, keys_values.get(), environment_id),
        )?;

        if rows_affected == 0 {
            return Err("Environment not found".into());
        }

        Ok(())
    }
}
