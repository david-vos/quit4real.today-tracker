-- Migration: Create Necessary Tables for the Application

-- Drop tables if they exist (optional, for safety when running tests)
-- DROP TABLE IF EXISTS game_failure_records;
-- DROP TABLE IF EXISTS platform_subscriptions;
-- DROP TABLE IF EXISTS games;
-- DROP TABLE IF EXISTS platforms;
-- DROP TABLE IF EXISTS users;
--

-- Create Users Table
CREATE TABLE IF NOT EXISTS users
(
    id   INTEGER PRIMARY KEY, -- User ID
    name TEXT NOT NULL        -- User's name
);

-- Create Platforms Table
CREATE TABLE IF NOT EXISTS platforms
(
    id   TEXT PRIMARY KEY,    -- Platform ID (e.g., 'steam')
    name TEXT UNIQUE NOT NULL -- Platform name (e.g., 'Steam')
);

-- Create User Platform Subscriptions Table
CREATE TABLE IF NOT EXISTS platform_subscriptions
(
    id               INTEGER PRIMARY KEY, -- Auto-incrementing ID
    display_name     TEXT,                -- User ID (foreign key)
    platform_id      TEXT,                -- Platform ID (foreign key)
    platform_game_id TEXT,                -- Game ID (foreign key)
    platform_user_id TEXT,                -- User's ID on the specific platform
    played_amount    INT DEFAULT 0,       -- Amount of time played
    FOREIGN KEY (platform_game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE
);

-- Create Games Table
CREATE TABLE IF NOT EXISTS games
(
    id          TEXT PRIMARY KEY, -- Game ID (e.g., Steam App ID)
    name        TEXT NOT NULL,    -- Game name
    platform_id TEXT,             -- Platform ID (foreign key)
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE
);

-- Create Game Failure Records Table
CREATE TABLE IF NOT EXISTS game_failure_records
(
    id               INTEGER PRIMARY KEY, -- Auto-incrementing ID
    display_name     TEXT,
    platform_id      TEXT,                -- Game ID (foreign key)
    platform_game_id TEXT,
    platform_user_id TEXT,
    duration_minutes INT,                 -- Duration of the failure in minutes
    reason           TEXT,                -- Reason for the failure
    timestamp        DATETIME,            -- When the failure was recorded
    FOREIGN KEY (platform_id) references platforms (id) ON DELETE CASCADE,
    FOREIGN KEY (platform_game_id) REFERENCES games (id) ON DELETE CASCADE
);

-- Migration completed
