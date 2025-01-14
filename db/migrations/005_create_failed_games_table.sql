CREATE TABLE IF NOT EXISTS failed_games (
                                            id INTEGER PRIMARY KEY AUTOINCREMENT,
                                            steam_id TEXT NOT NULL,
                                            game_id TEXT NOT NULL,
                                            failed_time DATETIME NOT NULL,
                                            FOREIGN KEY (steam_id) REFERENCES users(steam_id),
                                            UNIQUE (steam_id, game_id, failed_time)
);