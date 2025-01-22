-- Migration: Create Necessary Tables for the Application

-- Drop tables if they exist (optional, for safety)
DROP TABLE IF EXISTS game_failure_records;
DROP TABLE IF EXISTS user_platform_subscriptions;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS platforms;
DROP TABLE IF EXISTS users;

-- Create Users Table
CREATE TABLE users
(
    id   TEXT PRIMARY KEY, -- User ID (Steam ID)
    name TEXT NOT NULL     -- User's name
);

-- Create Platforms Table
CREATE TABLE platforms
(
    id   TEXT PRIMARY KEY,    -- Platform ID (e.g., 'steam')
    name TEXT UNIQUE NOT NULL -- Platform name (e.g., 'Steam')
);

-- Insert Steam Platform
INSERT INTO platforms (id, name)
VALUES ('steam', 'Steam');

-- Create User Platform Subscriptions Table
CREATE TABLE user_platform_subscriptions
(
    id               INTEGER PRIMARY KEY, -- Auto-incrementing ID
    user_id          TEXT,                -- User ID (foreign key)
    platform_id      TEXT,                -- Platform ID (foreign key)
    platform_game_id TEXT,                -- Game ID (foreign key)
    platform_user_id TEXT,                -- User's ID on the specific platform
    played_amount    INT DEFAULT 0,       -- Amount of time played
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (platform_game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE
);

-- Create Games Table
CREATE TABLE games
(
    id          TEXT PRIMARY KEY, -- Game ID (e.g., Steam App ID)
    name        TEXT NOT NULL,    -- Game name
    platform_id TEXT,             -- Platform ID (foreign key)
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE
);

-- Create Game Failure Records Table
CREATE TABLE game_failure_records
(
    id               INTEGER PRIMARY KEY, -- Auto-incrementing ID
    user_id          TEXT,                -- User ID (foreign key)
    game_id          TEXT,                -- Game ID (foreign key)
    duration_minutes INT,                 -- Duration of the failure in minutes
    reason           TEXT,                -- Reason for the failure
    timestamp        DATETIME,            -- When the failure was recorded
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
);

-- Migration completed
