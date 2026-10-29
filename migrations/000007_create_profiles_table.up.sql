CREATE TABLE profiles
(
    user_id    UUID PRIMARY KEY REFERENCES users (id) ON DELETE CASCADE,
    username   TEXT,
    first_name TEXT,
    last_name  TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

INSERT INTO profiles (user_id, username, created_at, updated_at)
SELECT id, username, created_at, updated_at
FROM users;

ALTER TABLE users
    DROP COLUMN username;