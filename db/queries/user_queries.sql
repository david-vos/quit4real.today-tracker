-- Get all users
SELECT id, name, email FROM users;

-- Insert a user
INSERT INTO users (name, email) VALUES (?, ?);
