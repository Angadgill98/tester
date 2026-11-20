use std::{collections::HashMap, sync::Arc};

use rusqlite::Connection;
use serde_json::Value;
use tokio_tungstenite::{
    connect_async,
    tungstenite::Message,
    WebSocketStream,
    MaybeTlsStream,
};

use tokio::net::TcpStream;
use futures_util::{SinkExt, StreamExt};

use crate::engine::{Action, Application, Sequence};

pub type WsConnection = WebSocketStream<MaybeTlsStream<TcpStream>>;

pub struct SequenceState {
    pub state: HashMap<i64, Vec<ActionResponse>>,
    pub order: Vec<i64>,
    pub custom_data: HashMap<String, Value>,
}

pub struct ActionResponse {
    pub success: bool,
    pub message: String,
    pub body: String,
    pub status_code: u16,
    pub status: String,
    pub headers: HashMap<String, Vec<String>>,
    pub cookies: Vec<Cookie>,
    pub data: HashMap<String, Value>,
}

pub struct Cookie {
    pub name: String,
    pub value: String,
}



type TokioRWLock<T> = tokio::sync::RwLock<T>;

pub struct Tester {
    pub sequence_id:String,
    pub execution_id:String,
    pub app_id:String,
    
    pub client: reqwest::Client,
    pub sequence_state: Arc<tokio::sync::Mutex<SequenceState>>,
    pub sub_sequence_context: HashMap<String, SubSequenceContext>,
    pub websocket_manager: Arc<TokioRWLock<HashMap<String, Arc<TokioRWLock<WsConnection>>>>>,
}

struct TesterManager {
    db:tokio::sync::Mutex<Connection>,
    sequence_obj:Sequence::sequence::SequenceObj,
    action_obj:Action::action::ActionObj,
    application_onj:Application::applciation::ApplicationObj,

    application_execution_context: Arc<TokioRWLock<HashMap<String, TokioRWLock<HashMap<String,ExecutionContext>>>>>,//appid and te exection id and the exeution cotext 
    sequence_execution_context: Arc<TokioRWLock<HashMap<String, TokioRWLock<Option<Tester>>>>>,//excetiion id adn the teh sequnce being executed 


    cancellation_tokens:HashMap<String,tokio_util::sync::CancellationToken>,

    execution_queue_sender: tokio::sync::mpsc::Sender<>,
    execution_queue_rec: tokio::sync::mpsc::Receiver<>,

    sequence_queue_sender: tokio::sync::mpsc::Sender<>,
    sequence_queue_rec: tokio::sync::mpsc::Receiver<>,

    sub_sequence_queue_sender: tokio::sync::mpsc::Sender<>,
    sub_sequence_queue_rec: tokio::sync::mpsc::Receiver<>,

    action_queue_sender: tokio::sync::mpsc::Sender<>,
    action_queue_rec: tokio::sync::mpsc::Receiver<>,


}

struct ExecutionContext{
    execution_id:String,
    app_id:String,
    sequences:Vec<SequenceContext>,
    total_sequences:u64,
    sequences_executed:u64,
    pub cancellation_token: tokio_util::sync::CancellationToken,

}

struct SequenceContext{
    sequence_id:String,


    sub_sequences:Vec<SubSequenceContext>,

    total_subsequnce:u64,
    sub_completed:u64,
    pub cancellation_token: tokio_util::sync::CancellationToken,

    completed:bool
}

struct SubSequenceContext {
    action_context: HashMap<String, ActionContext>,
    total_actions:u64,
    actions_completed:u64,
    pub cancellation_token: tokio_util::sync::CancellationToken,

}

pub struct ActionContext {
    pub req_context_id: String,
    pub action_id: String,
    pub action_type: String,
    pub cancellation_token: tokio_util::sync::CancellationToken,
}


impl TesterManager{
    
    pub fn ExecuteSequence(&self,app_id:String,sequence_id:String){

    }

    async fn HandleExecutionContext(&self, app_id: String, sequence_id: String)->Result<(),Box<dyn std::error::Error>> {
        let cancellation_token = tokio_util::sync::CancellationToken::new();
        let execution_id = uuid::Uuid::new_v4().to_string();
    
        let exe_ctx = ExecutionContext {
            app_id,
            execution_id,
            cancellation_token,
        };
    
        let db = self.db.lock().await;
        
        let mut sequence = self.sequence_obj.get_sequence(&db, sequence_id.as_str())?;
        
        let mut sequences = Vec::new();
        
        sequences.push(sequence.clone());
        
        while let Some(next_sequence_id) = sequence.next.clone() {
            sequence = self.sequence_obj.get_sequence(&db, next_sequence_id.as_str())?;
            sequences.push(sequence.clone());
        }


        let total_sequences=sequences.len();

        
        let actions=Vec::new();

        sequence.
        




        let exe_ctx = ExecutionContext {
            app_id,
            execution_id,
            cancellation_token,
            total_sequences:total_sequences as u64,
        };
        Ok(())
    }

    fn HandleAppContext(&mut self, app_id: String, sequence_id: String) {

    }

    fn HandleSequenceContext(&mut self, sequence_id: String, tester: Option<Tester>) {

    }


}

