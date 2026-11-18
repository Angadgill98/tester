use std::{collections::HashMap, sync::Arc};

use serde_json::Value;
use tokio_tungstenite::{
    connect_async,
    tungstenite::Message,
    WebSocketStream,
    MaybeTlsStream,
};

use tokio::net::TcpStream;
use futures_util::{SinkExt, StreamExt};

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
    pub client: reqwest::Client,
    pub sequence_state: Arc<tokio::sync::Mutex<SequenceState>>,
    pub sub_sequence_context: HashMap<String, SubSequenceContext>,
    pub websocket_manager: Arc<TokioRWLock<HashMap<String, Arc<TokioRWLock<WsConnection>>>>>,
}

struct TesterManager {
    application_execution_context: Arc<TokioRWLock<HashMap<String, TokioRWLock<Vec<String>>>>>,
    sequence_execution_context: Arc<TokioRWLock<HashMap<String, TokioRWLock<Option<Tester>>>>>,
}

struct SubSequenceContext {
    action_context: HashMap<String, ActionContext>,
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

    fn HandleAppContext(&mut self, app_id: String, sequence_id: String) {

    }

    fn HandleSequenceContext(&mut self, sequence_id: String, tester: Option<Tester>) {

    }


}