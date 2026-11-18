use rusqlite::Connection;



pub struct ActionObj{
    
}

impl ActionObj {
    pub fn create_action_obj() -> Self {
        Self {
        }
    }

    pub fn register_action(&self,db: &Connection,action_type:&String,action_id:&String)->Result<(),rusqlite::Error>{
        let query:String=format!("insert into actions (action_id,action_type) values(?1,?2)");
        db.execute(&query, (action_id,action_type))?;

        Ok(())
    }

    pub fn register_actions(&self,db: &Connection, action_type: Vec<String>, action_id: Vec<String>) -> Result<(), rusqlite::Error> {
        if action_type.len() != action_id.len() {
            return Err(rusqlite::Error::InvalidParameterCount(action_type.len(), action_id.len()));
        }
    
        if action_id.is_empty() {
            return Ok(());
        }
    
        let values = (0..action_id.len())
            .map(|i| format!("(?{}, ?{})", i * 2 + 1, i * 2 + 2))
            .collect::<Vec<String>>()
            .join(", ");
    
        let query = format!("INSERT INTO actions (action_id, action_type) VALUES {}", values);
    
        let params: Vec<&String> = action_id.iter()
            .zip(action_type.iter())
            .flat_map(|(id, action)| [id, action])
            .collect();
    
        db.execute(&query, rusqlite::params_from_iter(params))?;
    
        Ok(())
    }

    pub fn get_action(&self, db: &Connection, action_id: String) -> Result<Action, rusqlite::Error> {
        let query = "
            SELECT action_id, action_type
            FROM actions
            WHERE action_id = ?1
        ";
    
        let action = db.query_row(query, [&action_id], |row| {
            Ok(Action {
                action_id: row.get(0)?,
                action_type: row.get(1)?,
            })
        })?;
    
        Ok(action)
    }
    
    pub fn get_actions(&self, db: &Connection, action_ids: Vec<String>) -> Result<Vec<Action>, rusqlite::Error> {
        if action_ids.is_empty() {
            return Ok(Vec::new());
        }
    
        let placeholders = (1..=action_ids.len())
            .map(|i| format!("?{}", i))
            .collect::<Vec<String>>()
            .join(", ");
    
        let query = format!(
            "SELECT action_id, action_type FROM actions WHERE action_id IN ({})",
            placeholders
        );
    
        let mut stmt = db.prepare(&query)?;
    
        let actions = stmt.query_map(rusqlite::params_from_iter(action_ids.iter()), |row| {
            Ok(Action {
                action_id: row.get(0)?,
                action_type: row.get(1)?,
            })
        })?;
    
        actions.collect()
    }
    pub fn delete_action(&self,db: &Connection, action_id: String) -> Result<(), rusqlite::Error> {
        let query = "DELETE FROM actions WHERE action_id = ?1";
    
        db.execute(query, [&action_id])?;
    
        Ok(())
    }
    
}


pub struct Action{
    pub action_id:String,
    pub action_type:String
}

