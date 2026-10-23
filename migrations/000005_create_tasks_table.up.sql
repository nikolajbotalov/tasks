CREATE TABLE tasks
(
    id          UUID PRIMARY KEY,
    name        TEXT        NOT NULL,
    description TEXT        NOT NULL,
    author_id   UUID        NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    assignee_id UUID        REFERENCES users (id) ON DELETE SET NULL,
    status      TEXT,
    priority    TEXT,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX tasks_author_id_idx ON tasks (author_id);
CREATE INDEX tasks_assignee_id_idx ON tasks (assignee_id);