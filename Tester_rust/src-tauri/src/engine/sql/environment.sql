CREATE TABLE IF NOT EXISTS environments (
    environment_id TEXT PRIMARY KEY,
    environment_name TEXT NOT NULL,
    keys_values TEXT NOT NULL DEFAULT '{}'
);