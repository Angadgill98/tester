use std::collections::HashMap;
use serde::{Deserialize, Serialize};

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct SequenceState {
    pub state: HashMap<String, Vec<ActionResponse>>,
    pub order: Vec<String>,
    pub custom_data: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ActionResponse {
    pub success: bool,
    pub message: String,
    pub body: String,
    pub status_code: u16,
    pub status: String,
    pub headers: HashMap<String, Vec<String>>,
    pub cookies: Vec<Cookie>,
    pub data: HashMap<String, serde_json::Value>,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct Cookie {
    pub name: String,
    pub value: String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ActionContext {
    pub req_context_id: i32,
    pub action_id: i32,
    pub action_type: String,
}

pub struct RequestExecutionContext {
    pub ctx: tokio_util::sync::CancellationToken,
}