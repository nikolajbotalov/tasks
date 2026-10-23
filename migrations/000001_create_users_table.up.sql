CREATE TABLE users
(
    id         UUID PRIMARY KEY,
    email      TEXT UNIQUE NOT NULL,
    username   TEXT        NOT NULL,
    password   TEXT        NOT NULL,
    role       TEXT        NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);