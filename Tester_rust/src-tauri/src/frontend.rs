
use std::{collections::HashMap};

use serde::{Deserialize, Serialize};
use serde_json::value::RawValue;

use crate::{AppState, engine::{Action::{self, action}, Application, Sequence}};

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct ActionInfo {
    pub id: String,
    pub name: String,
    #[serde(rename = "Type")]
    pub action_type: String,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct SequenceInfo {
    pub id: String,
    pub name: String,
}

#[derive(Debug, Serialize, Deserialize)]
#[serde(rename_all = "PascalCase")]
pub struct ChildApplicationInfo {
    pub id: String,
    pub name: String,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Response {
    pub success: bool,
    pub data: serde_json::Value,
    pub error: String,
}

pub fn convert_actions(actions: Vec<action::Action>) -> Vec<ActionInfo> {
    actions.into_iter().map(|action| ActionInfo {
        id: action.action_id.to_string(),
        name: action.action_type.clone(),
        action_type: action.action_type,
    }).collect()
}

pub fn convert_sequences(sequences: Vec<Sequence::sequence::Sequence>) -> Vec<SequenceInfo> {
    sequences.into_iter().map(|sequence| SequenceInfo {
        id: sequence.sequence_id.to_string(),
        name: sequence.sequence_name,
    }).collect()
}

pub fn convert_children(children: Vec<Application::applciation::Application>) -> Vec<ChildApplicationInfo> {
    children.into_iter().map(|child| ChildApplicationInfo {
        id: child.app_id.to_string(),
        name: child.app_name,
    }).collect()
}


#[tauri::command]
pub fn get_applications(state: tauri::State<'_, AppState>, app_ids: Vec<String>) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => {
            return Response {
                success: false,
                data: serde_json::Value::Null,
                error: e.to_string(),
            };
        }
    };

    let applications = match state.application_obj.get_applications_by_ids(&db, app_ids) {
        Ok(apps) => apps,
        Err(e) => {
            return Response {
                success: false,
                data: serde_json::Value::Null,
                error: e.to_string(),
            };
        }
    };

    let mut results = Vec::new();

    for app in applications {
        let actions = match state.action_obj.get_actions(&db, app.actions.clone()) {
            Ok(actions) => actions,
            Err(e) => {
                return Response {
                    success: false,
                    data: serde_json::Value::Null,
                    error: e.to_string(),
                };
            }
        };

        let mut action_info = Vec::new();

        for action in actions {
            let action_name = match action.action_type.as_str() {
                "http" => {
                    match state.action_obj.get_http_action(&db, action.action_id.as_str()) {
                        Ok(http_action) => http_action.action_name,
                        Err(e) => {
                            return Response {
                                success: false,
                                data: serde_json::Value::Null,
                                error: e.to_string(),
                            };
                        }
                    }
                }

                "ws" => {
                    match state.action_obj.get_ws_action(&db, action.action_id.as_str()) {
                        Ok(ws_action) => ws_action.action_name,
                        Err(e) => {
                            return Response {
                                success: false,
                                data: serde_json::Value::Null,
                                error: e.to_string(),
                            };
                        }
                    }
                }

                _ => action.action_type.clone(),
            };

            action_info.push(ActionInfo {
                id: action.action_id,
                name: action_name,
                action_type: action.action_type,
            });
        }

        let sequences = match state.sequnce_obj.get_sequences(&db, app.sequences.clone()) {
            Ok(sequences) => sequences,
            Err(e) => {
                return Response {
                    success: false,
                    data: serde_json::Value::Null,
                    error: e.to_string(),
                };
            }
        };

        let children = match state.application_obj.get_applications_by_ids(&db, app.children.clone()) {
            Ok(children) => children,
            Err(e) => {
                return Response {
                    success: false,
                    data: serde_json::Value::Null,
                    error: e.to_string(),
                };
            }
        };

        results.push(serde_json::json!({
            "app": app,
            "actions": action_info,
            "sequences": convert_sequences(sequences),
            "children": convert_children(children)
        }));
    }

    Response {
        success: true,
        data: serde_json::json!(results),
        error: String::new(),
    }
}

fn error_response(error: impl ToString) -> Response {
    Response {
        success: false,
        data: serde_json::Value::Null,
        error: error.to_string(),
    }
}

#[tauri::command]
pub fn create_application(state: tauri::State<'_, AppState>, app_name: String, parent_app_id: String) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => return error_response(e),
    };

    let app = match state.application_obj.create_application(&db, &app_name, parent_app_id.clone()) {
        Ok(app) => app,
        Err(e) => return error_response(e),
    };

    if let Err(e) = state.application_obj.register_child(&db, parent_app_id, app.app_id.clone()) {
        return error_response(e);
    }

    Response {
        success: true,
        data: serde_json::json!(app),
        error: String::new(),
    }
}

#[tauri::command]
pub fn create_http_action(state: tauri::State<'_, AppState>, action_name: String, action_type: String, app_id: String) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => return error_response(e),
    };

    let default_opts = Action::HttpAction::HttpActionOpts {
        method: "GET".to_string(),
        url: String::new(),
        headers: HashMap::new(),
        query: HashMap::new(),
        body: RawValue::from_string("{}".to_string()).unwrap(),
        cookies: HashMap::new(),
    };

    let environment_id = "0";

    let action_id = match state.action_obj.create_http_action(&db, &action_name, &app_id, &default_opts, environment_id) {
        Ok(id) => id,
        Err(e) => return error_response(e),
    };

    if let Err(e) = db.execute(
        "INSERT INTO actions (action_id, action_type) VALUES (?1, ?2)",
        (&action_id, &action_type),
    ) {
        return error_response(e);
    }

    if let Err(e) = state.application_obj.register_action(&db, app_id.clone(), action_id.clone()) {
        return error_response(e);
    }

    Response {
        success: true,
        data: serde_json::json!(action_id),
        error: String::new(),
    }
}

#[tauri::command]
pub fn create_ws_action(state: tauri::State<'_, AppState>, action_name: String, action_type: String, app_id: String) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => return error_response(e),
    };

    let default_opts = Action::WsAction::WSOpts {
        url: String::new(),
        headers: HashMap::new(),
        query_params: HashMap::new(),
        subprotocols: Vec::new(),
        messages: Vec::new(),
        handshake_timeout:5,
        write_timeout: 10,
    };

    let environment_id = "0";

    let action_id = match state.action_obj.create_ws_action(&db, &app_id, &action_name, "", &default_opts, environment_id) {
        Ok(id) => id,
        Err(e) => return error_response(e),
    };

    if let Err(e) = db.execute(
        "INSERT INTO actions (action_id, action_type) VALUES (?1, ?2)",
        (&action_id, &action_type),
    ) {
        return error_response(e);
    }

    if let Err(e) = state.application_obj.register_action(&db, app_id.clone(), action_id.clone()) {
        return error_response(e);
    }

    Response {
        success: true,
        data: serde_json::json!(action_id),
        error: String::new(),
    }
}

#[tauri::command]
pub fn get_http_action(state: tauri::State<'_, AppState>, action_id: String) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => return error_response(e),
    };

    let action = match state.action_obj.get_http_action(&db, &action_id) {
        Ok(action) => action,
        Err(e) => return error_response(e),
    };

    Response {
        success: true,
        data: serde_json::json!(action),
        error: String::new(),
    }
}

#[tauri::command]
pub fn get_ws_action(state: tauri::State<'_, AppState>, action_id: String) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => return error_response(e),
    };

    let action = match state.action_obj.get_ws_action(&db, &action_id) {
        Ok(action) => action,
        Err(e) => return error_response(e),
    };

    Response {
        success: true,
        data: serde_json::json!(action),
        error: String::new(),
    }
}

#[tauri::command]
pub fn get_sequence(state: tauri::State<'_, AppState>, sequence_id: String) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => return error_response(e),
    };

    let sequence = match state.sequnce_obj.get_sequence(&db, &sequence_id) {
        Ok(sequence) => sequence,
        Err(e) => return error_response(e),
    };

    Response {
        success: true,
        data: serde_json::json!(sequence),
        error: String::new(),
    }
}

#[tauri::command]
pub fn create_sequence(state: tauri::State<'_, AppState>, sequence_name: String, app_id: String) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => return error_response(e),
    };

    let sequence = match state.sequnce_obj.create_sequence(&db, &sequence_name, &app_id) {
        Ok(sequence) => sequence,
        Err(e) => return error_response(e),
    };

    if let Err(e) = state.application_obj.register_sequence(&db, app_id.clone(), sequence.sequence_id.clone()) {
        return error_response(e);
    }

    Response {
        success: true,
        data: serde_json::json!(sequence),
        error: String::new(),
    }
}

#[tauri::command]
pub fn save_http_action(state: tauri::State<'_, AppState>, action: Action::HttpAction::HttpAction) -> Response {
    let db = match state.db.lock() {
        Ok(db) => db,
        Err(e) => return Response {
            success: false,
            data: serde_json::Value::Null,
            error: e.to_string(),
        },
    };

    match state.action_obj.update_http_action(&db, action.action_id.as_str(), action.action_name.as_str(), action.app_id.as_str(), &action.opts, action.environment_id.as_str()) {
        Ok(_) => Response {
            success: true,
            data: serde_json::Value::Null,
            error: String::new(),
        },
        Err(e) => Response {
            success: false,
            data: serde_json::Value::Null,
            error: e.to_string(),
        },
    }
}