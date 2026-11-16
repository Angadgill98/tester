
CREATE TABLE IF NOT EXISTS sequences (
    sequence_id TEXT PRIMARY KEY,
    app_id TEXT NOT NULL,
    sequence_name TEXT NOT NULL,
    action_ids TEXT NOT NULL DEFAULT '[]'
);