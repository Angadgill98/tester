CREATE TABLE IF NOT EXISTS sequences (
    sequence_id INTEGER PRIMARY KEY AUTOINCREMENT,
    app_id INTEGER NOT NULL,
    sequence_name TEXT NOT NULL,
    action_ids TEXT NOT NULL,

    FOREIGN KEY (app_id) REFERENCES apps(app_id)
);