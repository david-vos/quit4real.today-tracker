CREATE TABLE IF NOT EXISTS tracker (
                                       user_id INTEGER NOT NULL,
                                       game_id INTEGER NOT NULL,
                                       platform_id VARCHAR NOT NULL,
                                       day DATETIME NOT NULL,
                                       time_played INTEGER NOT NULL,
                                       new_total_time_played INTEGER NOT NULL,
                                       amount_of_logins INTEGER NOT NULL,
                                       PRIMARY KEY (user_id, game_id, platform_id, day)
);
