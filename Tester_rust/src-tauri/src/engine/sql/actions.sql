-- 1. Main actions table
CREATE TABLE IF NOT EXISTS actions (
    action_id INTEGER PRIMARY KEY,
    action_type TEXT NOT NULL
);

-- 2. HTTP actions table
CREATE TABLE IF NOT EXISTS http_actions (
    action_id INTEGER PRIMARY KEY,
    action_name TEXT NOT NULL,
    app_id TEXT NOT NULL,
    environment_id TEXT,
    opts TEXT NOT NULL DEFAULT '{}'
);

-- 3. WebSocket actions table
CREATE TABLE IF NOT EXISTS ws_actions (
    action_id INTEGER PRIMARY KEY,
    action_name TEXT NOT NULL,
    app_id TEXT NOT NULL,
    environment_id TEXT,
    opts TEXT NOT NULL DEFAULT '{}'
);