use rusqlite::{params_from_iter, Connection};
use serde::{Deserialize, Serialize};
use std::error::Error;

pub struct SequenceObj;

impl SequenceObj {
    pub fn new_sequence_obj() -> Self {
        Self
    }
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct Sequence {
    pub sequence_id: String,
    pub app_id: String,
    pub sequence_name: String,
    pub action_ids: Vec<String>,
}

impl SequenceObj {
    pub fn create_sequence(&self,db: &Connection, sequence_name: &str, app_id: &str) -> Result<Sequence, Box<dyn Error>> {
        let sequence_id = uuid::Uuid::new_v4().to_string();
    
        let mut sequence = Sequence {
            sequence_id,
            app_id: app_id.to_string(),
            sequence_name: sequence_name.to_string(),
            action_ids: Vec::new(),
        };
    
        let action_ids = serde_json::to_string(&sequence.action_ids)?;
    
        let query = "
            INSERT INTO sequences (sequence_id, app_id, sequence_name, action_ids)
            VALUES (?1, ?2, ?3, ?4)
        ";
    
        db.execute(query, (&sequence.sequence_id, &sequence.app_id, &sequence.sequence_name, action_ids))?;
    
        Ok(sequence)
    }
    
    pub fn create_temporary_sequence(sequence_name: &str, app_id: &str) -> Sequence {
        Sequence {
            sequence_id: String::new(),
            app_id: app_id.to_string(),
            sequence_name: sequence_name.to_string(),
            action_ids: Vec::new(),
        }
    }

    pub fn get_sequence(&self,db: &Connection, sequence_id: &str) -> Result<Sequence, Box<dyn Error>> {
        let query = "
            SELECT sequence_id, app_id, sequence_name, action_ids
            FROM sequences
            WHERE sequence_id = ?1
        ";

        let (sequence_id, app_id, sequence_name, action_ids_json) =
            db.query_row(query, [sequence_id], |row| {
                Ok((
                    row.get::<_, String>(0)?,
                    row.get::<_, String>(1)?,
                    row.get::<_, String>(2)?,
                    row.get::<_, String>(3)?,
                ))
            })?;

        let action_ids = serde_json::from_str(&action_ids_json)?;

        Ok(Sequence {
            sequence_id,
            app_id,
            sequence_name,
            action_ids,
        })
    }

    pub fn get_sequences(&self,db: &Connection, sequence_ids: Vec<String>) -> Result<Vec<Sequence>, Box<dyn Error>> {
        if sequence_ids.is_empty() {
            return Ok(Vec::new());
        }

        let placeholders = (1..=sequence_ids.len())
            .map(|i| format!("?{}", i))
            .collect::<Vec<String>>()
            .join(", ");

        let query = format!(
            "SELECT sequence_id, app_id, sequence_name, action_ids
             FROM sequences
             WHERE sequence_id IN ({})",
            placeholders
        );

        let mut stmt = db.prepare(&query)?;

        let rows = stmt.query_map(params_from_iter(sequence_ids.iter()), |row| {
            Ok((
                row.get::<_, String>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, String>(2)?,
                row.get::<_, String>(3)?,
            ))
        })?;

        let mut sequences = Vec::new();

        for row in rows {
            let (sequence_id, app_id, sequence_name, action_ids_json) = row?;
            let action_ids = serde_json::from_str(&action_ids_json)?;

            sequences.push(Sequence {
                sequence_id,
                app_id,
                sequence_name,
                action_ids,
            });
        }

        Ok(sequences)
    }

    pub fn update_sequence(&self,db: &Connection, sequence_id: &str, sequence_name: &str, action_ids: &[String]) -> Result<(), Box<dyn Error>> {
        let action_ids_json = serde_json::to_string(action_ids)?;

        let query = "
            UPDATE sequences
            SET sequence_name = ?1,
                action_ids = ?2
            WHERE sequence_id = ?3
        ";

        db.execute(query, (sequence_name, action_ids_json, sequence_id))?;

        Ok(())
    }

    pub fn delete_sequence(&self,db: &Connection, sequence_id: &str) -> Result<(), Box<dyn Error>> {
        let query = "DELETE FROM sequences WHERE sequence_id = ?1";
        db.execute(query, [sequence_id])?;

        Ok(())
    }
}