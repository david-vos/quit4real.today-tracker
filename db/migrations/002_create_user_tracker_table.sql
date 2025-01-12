CREATE TABLE IF NOT EXISTS user_tracker (
                                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                                            steam_id TEXT NOT NULL,
                                            game_id TEXT NOT NULL,
                                            played_amount INT,
                                            FOREIGN KEY (steam_id) REFERENCES users(steam_id),
                                            UNIQUE (steam_id, game_id)
);
