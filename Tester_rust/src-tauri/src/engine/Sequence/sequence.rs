
use rusqlite::{params_from_iter, Connection};
use std::{collections::HashMap, error::Error};

use crate::engine::Action;

pub struct SequenceObj;

impl SequenceObj {
    pub fn new_sequence_obj() -> Self {
        Self
    }
}

pub struct Sequence {
    pub sequence_id: String,
    pub app_id: String,
    pub sequence_name: String,
    pub actions: Vec<Action::action::Action>,
    pub root_sequences: HashMap<String, SubSequence>,
    pub sub_sequences: HashMap<String, SubSequence>,
    pub next: Option<String>,
}
pub struct SubSequence {
    pub sub_id:String,
    pub order:String,
    pub actions: Vec<String>,
    pub next: Vec<String>,
}

impl SequenceObj {
    pub fn create_sequence(&self, db: &Connection, sequence_name: &str, app_id: &str, sub_sequences: &[u8], next: Option<&str>) -> Result<SequenceRepo, Box<dyn Error>> {
        let sequence = SequenceRepo {
            sequence_id: uuid::Uuid::new_v4().to_string(),
            app_id: app_id.to_string(),
            sequence_name: sequence_name.to_string(),
            action_ids: Vec::new(),
            root_sequences: Vec::new(),
            sub_sequences: sub_sequences.to_vec(),
            next: next.map(String::from),
        };

        let action_ids = serde_json::to_string(&sequence.action_ids)?;
        let root_sequences = serde_json::to_string(&sequence.root_sequences)?;

        let query = "
            INSERT INTO sequences (
                sequence_id, app_id, sequence_name,
                action_ids, root_sequences, sub_sequences, next
            )
            VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)
        ";

        db.execute(query, (
            &sequence.sequence_id,
            &sequence.app_id,
            &sequence.sequence_name,
            action_ids,
            root_sequences,
            &sequence.sub_sequences,
            &sequence.next,
        ))?;

        Ok(sequence)
    }

    pub fn create_temporary_sequence(sequence_name: &str, app_id: &str, sub_sequences: &[u8], next: Option<&str>) -> SequenceRepo {
        SequenceRepo {
            sequence_id: String::new(),
            app_id: app_id.to_string(),
            sequence_name: sequence_name.to_string(),
            action_ids: Vec::new(),
            root_sequences: Vec::new(),
            sub_sequences: sub_sequences.to_vec(),
            next: next.map(String::from),
        }
    }

    pub fn get_sequence(&self, db: &Connection, sequence_id: &str) -> Result<SequenceRepo, Box<dyn Error>> {
        let query = "
            SELECT sequence_id, app_id, sequence_name,
                   action_ids, root_sequences, sub_sequences, next
            FROM sequences
            WHERE sequence_id = ?1
        ";

        let sequence = db.query_row(query, [sequence_id], |row| {
            Ok((
                row.get::<_, String>(0)?,
                row.get::<_, String>(1)?,
                row.get::<_, String>(2)?,
                row.get::<_, String>(3)?,
                row.get::<_, String>(4)?,
                row.get::<_, Vec<u8>>(5)?,
                row.get::<_, Option<String>>(6)?,
            ))
        })?;

        Ok(SequenceRepo {
            sequence_id: sequence.0,
            app_id: sequence.1,
            sequence_name: sequence.2,
            action_ids: serde_json::from_str(&sequence.3)?,
            root_sequences: serde_json::from_str(&sequence.4)?,
            sub_sequences: sequence.5,
            next: sequence.6,
        })
    }

    pub fn get_sequences(&self, db: &Connection, sequence_ids: Vec<String>) -> Result<Vec<SequenceRepo>, Box<dyn Error>> {
        if sequence_ids.is_empty() {
            return Ok(Vec::new());
        }

        let placeholders = (1..=sequence_ids.len())
            .map(|i| format!("?{}", i))
            .collect::<Vec<String>>()
            .join(", ");

        let query = format!(
            "SELECT sequence_id, app_id, sequence_name,
                    action_ids, root_sequences, sub_sequences, next
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
                row.get::<_, String>(4)?,
                row.get::<_, Vec<u8>>(5)?,
                row.get::<_, Option<String>>(6)?,
            ))
        })?;

        let mut sequences = Vec::new();

        for row in rows {
            let (sequence_id, app_id, sequence_name, action_ids_json, root_sequences_json, sub_sequences, next) = row?;

            sequences.push(SequenceRepo {
                sequence_id,
                app_id,
                sequence_name,
                action_ids: serde_json::from_str(&action_ids_json)?,
                root_sequences: serde_json::from_str(&root_sequences_json)?,
                sub_sequences,
                next,
            });
        }

        Ok(sequences)
    }

    pub fn update_sequence(&self, db: &Connection, sequence_id: &str, sequence_name: &str, action_ids: &[String], root_sequences: &[String], sub_sequences: &[u8], next: Option<&str>) -> Result<(), Box<dyn Error>> {
        let action_ids_json = serde_json::to_string(action_ids)?;
        let root_sequences_json = serde_json::to_string(root_sequences)?;

        let query = "
            UPDATE sequences
            SET sequence_name = ?1,
                action_ids = ?2,
                root_sequences = ?3,
                sub_sequences = ?4,
                next = ?5
            WHERE sequence_id = ?6
        ";

        db.execute(query, (
            sequence_name,
            action_ids_json,
            root_sequences_json,
            sub_sequences,
            next,
            sequence_id,
        ))?;

        Ok(())
    }

    pub fn delete_sequence(&self, db: &Connection, sequence_id: &str) -> Result<(), Box<dyn Error>> {
        let query = "DELETE FROM sequences WHERE sequence_id = ?1";
        db.execute(query, [sequence_id])?;

        Ok(())
    }
}