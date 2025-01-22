-- Migration: Create Users, Platforms, Games, and Game Failure Records Tables

-- Drop tables if they exist (optional, for safety)
DROP TABLE IF EXISTS game_failure_records;
DROP TABLE IF EXISTS games;
DROP TABLE IF EXISTS user_platform_subscriptions;
DROP TABLE IF EXISTS platforms;
DROP TABLE IF EXISTS users;

-- Create Users Table
CREATE TABLE users
(
    id       TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL
);

-- Create Platforms Table
CREATE TABLE platforms
(
    id   TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL
);

-- Create User Platform Subscriptions Table
CREATE TABLE user_platform_subscriptions
(
    id          INTEGER PRIMARY KEY,
    user_id     TEXT,
    game_id     TEXT,
    platform_id TEXT,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE,
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE
);

-- Create Games Table
CREATE TABLE games
(
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    platform_id TEXT,
    FOREIGN KEY (platform_id) REFERENCES platforms (id) ON DELETE CASCADE
);

-- Create Game Failure Records Table
CREATE TABLE game_fails_records
(
    id               INTEGER PRIMARY KEY,
    user_id          TEXT,
    game_id          TEXT,
    duration_minutes INTEGER,
    reason           TEXT,
    timestamp        DATETIME,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (game_id) REFERENCES games (id) ON DELETE CASCADE
);

-- Migration completed
