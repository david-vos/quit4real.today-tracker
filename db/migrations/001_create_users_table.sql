CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     name TEXT NOT NULL,
                                     steam_id TEXT NOT NULL,
                                     api_key TEXT NOT NULL
);
