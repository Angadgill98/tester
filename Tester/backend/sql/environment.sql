CREATE TABLE IF NOT EXISTS environments (
    environment_id INTEGER PRIMARY KEY AUTOINCREMENT,
    environment_name TEXT NOT NULL,
    keys_values TEXT NOT NULL DEFAULT '{}'
);