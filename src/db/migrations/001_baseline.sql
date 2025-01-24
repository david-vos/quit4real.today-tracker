-- Migration: Create Necessary Tables for the Application

-- Drop tables if they exist (optional, for safety when running tests)
-- DROP TABLE IF EXISTS game_failure_records;
-- DROP TABLE IF EXISTS platform_subscriptions;
-- DROP TABLE IF EXISTS games;
-- DROP TABLE IF EXISTS platforms;
-- DROP TABLE IF EXISTS users;

-- Create Users Table
CREATE TABLE IF NOT EXISTS users
(
    id   INTEGER PRIMARY KEY, -- User ID
    name TEXT UNIQUE NOT NULL  -- User's name (must be unique)
);

-- Create Platforms Table
CREATE TABLE IF NOT EXISTS platforms
(
    id   TEXT PRIMARY KEY,    -- Platform ID (e.g., 'steam')
    name TEXT UNIQUE NOT NULL  -- Platform name (e.g., 'Steam')
);

-- Create Games Table
CREATE TABLE IF NOT EXISTS games
(
    id          TEXT PRIMARY KEY, -- Game ID (e.g., Steam App ID)
    name        TEXT NOT NULL,    -- Game name
    platform_id TEXT NOT NULL,    -- Platform ID (foreign key, cannot be null)
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE,
    UNIQUE (id, platform_id)      -- Ensure only 1 game ID per platform
);

-- Create User Platform Subscriptions Table
CREATE TABLE IF NOT EXISTS platform_subscriptions
(
    id               INTEGER PRIMARY KEY, -- Auto-incrementing ID
    display_name     TEXT NOT NULL,       -- User's display name (cannot be null)
    platform_id      TEXT NOT NULL,       -- Platform ID (foreign key, cannot be null)
    platform_game_id TEXT NOT NULL,       -- Game ID (foreign key, cannot be null)
    platform_user_id TEXT NOT NULL,       -- User's ID on the specific platform (cannot be null)
    played_amount    INT DEFAULT 0,       -- Amount of time played
    FOREIGN KEY (platform_game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE,
    UNIQUE (platform_id, platform_game_id, platform_user_id) -- Ensure unique game per platform
);

-- Create Game Failure Records Table
CREATE TABLE IF NOT EXISTS game_failure_records
(
    id               INTEGER PRIMARY KEY, -- Auto-incrementing ID
    display_name     TEXT NOT NULL,       -- User's display name (cannot be null)
    platform_id      TEXT NOT NULL,       -- Platform ID (foreign key, cannot be null)
    platform_game_id TEXT NOT NULL,       -- Game ID (foreign key, cannot be null)
    platform_user_id TEXT NOT NULL,       -- User's ID on the specific platform (cannot be null)
    duration_minutes INT NOT NULL,        -- Duration of the failure in minutes (cannot be null)
    reason           TEXT NOT NULL,       -- Reason for the failure (cannot be null)
    timestamp        DATETIME NOT NULL,   -- When the failure was recorded (cannot be null)
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE,
    FOREIGN KEY (platform_game_id) REFERENCES games (id) ON DELETE CASCADE
);

-- Migration completed
