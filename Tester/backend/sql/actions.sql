-- 1. Main actions table
CREATE TABLE IF NOT EXISTS actions (
    action_id INTEGER PRIMARY KEY AUTOINCREMENT,
    action_type TEXT NOT NULL
);

-- 2. HTTP actions table
CREATE TABLE IF NOT EXISTS http_actions (
    action_id INTEGER PRIMARY KEY AUTOINCREMENT,
    action_name TEXT NOT NULL,
    app_id INTEGER NOT NULL,
    environment_id INTEGER,
    FOREIGN KEY (action_id) REFERENCES actions(action_id)
);