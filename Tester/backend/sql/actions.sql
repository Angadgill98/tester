-- 1. Main actions table
CREATE TABLE IF NOT EXISTS actions (
    action_id INTEGER PRIMARY KEY AUTOINCREMENT,
    action_type TEXT NOT NULL
);

-- 2. HTTP actions table
CREATE TABLE IF NOT EXISTS http_actions (
    action_id INTEGER PRIMARY KEY ,
    action_name TEXT NOT NULL,
    app_id INTEGER NOT NULL,
    environment_id INTEGER,
    opts TEXT NOT NULL DEFAULT '{}',
    FOREIGN KEY (action_id) REFERENCES actions(action_id)
);

CREATE TABLE IF NOT EXISTS ws_actions (
    action_id INTEGER PRIMARY KEY ,
    action_name TEXT NOT NULL,
    app_id INTEGER NOT NULL,
    environment_id INTEGER,
    opts TEXT NOT NULL DEFAULT '{}',
    FOREIGN KEY (action_id) REFERENCES actions(action_id),
    FOREIGN KEY (app_id) REFERENCES applications(id),
    FOREIGN KEY (environment_id) REFERENCES environments(id)
);