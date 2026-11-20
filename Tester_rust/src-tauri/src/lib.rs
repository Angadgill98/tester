mod engine;
mod init;
mod frontend;

use rusqlite::Connection;
use std::sync::Mutex;
use tauri::Manager;

use engine::Action::action;
use crate::engine::{Application, Sequence};

pub struct AppState {
    pub db: Mutex<Connection>,
    pub action_obj: action::ActionObj,
    pub application_obj: Application::applciation::ApplicationObj,
    pub sequnce_obj: Sequence::sequence::SequenceObj,
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        .setup(|app| {
            if cfg!(debug_assertions) {
                app.handle().plugin(
                    tauri_plugin_log::Builder::default()
                        .level(log::LevelFilter::Info)
                        .build(),
                )?;
            }

            let db_path = std::path::PathBuf::from(env!("CARGO_MANIFEST_DIR"))
                .join("../testflow.db");
            let db = init::init_db(db_path)?;

            app.manage(AppState {
                db: Mutex::new(db),
                action_obj: action::ActionObj {},
                application_obj: Application::applciation::ApplicationObj::create_application_obj(),
                sequnce_obj: Sequence::sequence::SequenceObj::new_sequence_obj(),
            });

            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            frontend::create_application,
            frontend::create_http_action,
            frontend::create_ws_action,
            frontend::get_http_action,
            frontend::get_ws_action,
            frontend::get_sequence,
            frontend::create_sequence,
            frontend::get_applications,
            frontend::save_http_action
        ])
        .run(tauri::generate_context!())
        .expect("error while building tauri application");
}