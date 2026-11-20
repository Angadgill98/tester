CREATE TABLE IF NOT EXISTS applications (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    parent_app_id INTEGER,
    environment_id INTEGER,
    action_ids TEXT NOT NULL DEFAULT '[]',
    sequence_ids TEXT NOT NULL DEFAULT '[]',
    children TEXT NOT NULL DEFAULT '[]',
    FOREIGN KEY (parent_app_id) REFERENCES applications(id) ON DELETE CASCADE,
    FOREIGN KEY (environment_id) REFERENCES environments(environment_id)
);