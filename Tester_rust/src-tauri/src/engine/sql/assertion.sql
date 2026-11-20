CREATE TABLE IF NOT EXISTS assertions (
    assertion_id INTEGER PRIMARY KEY AUTOINCREMENT,
    action_id INTEGER NOT NULL,
    status_code INTEGER NOT NULL,
    assertion TEXT NOT NULL DEFAULT '[]'
        CHECK (json_valid(assertion) AND json_type(assertion) = 'array'),
    FOREIGN KEY (action_id) REFERENCES actions(action_id) ON DELETE CASCADE
);