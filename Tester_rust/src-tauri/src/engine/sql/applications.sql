CREATE TABLE IF NOT EXISTS applications (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    parent_app_id TEXT,
    environment_id TEXT,
    action_ids TEXT NOT NULL DEFAULT '[]',
    sequence_ids TEXT NOT NULL DEFAULT '[]',
    children TEXT NOT NULL DEFAULT '[]'
);
