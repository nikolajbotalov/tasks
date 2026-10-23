INSERT INTO roles(name)
VALUES ('admin'),
       ('user'),
       ('viewer')
ON CONFLICT (name) DO NOTHING;

INSERT INTO role_permissions (role, permission)
VALUES ('admin', 'task:read'),
       ('admin', 'task:create'),
       ('admin', 'task:update'),
       ('admin', 'task:delete'),
       ('admin', 'comment:read'),
       ('admin', 'comment:create'),
       ('admin', 'user:read'),
       ('admin', 'user:update'),

       ('user', 'task:read'),
       ('user', 'task:create'),
       ('user', 'task:update'),
       ('user', 'task:delete'),
       ('user', 'comment:read'),
       ('user', 'comment:create'),

       ('viewer', 'task:read'),
       ('viewer', 'comment:read')
ON CONFLICT DO NOTHING;