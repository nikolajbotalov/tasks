INSERT INTO role_permissions (role, permission)
VALUES ('admin', 'profile:read'),
       ('admin', 'profile:update'),
       ('user', 'profile:read'),
       ('user', 'profile:update'),
       ('viewer', 'profile:read'),
       ('viewer', 'profile:update')
ON CONFLICT DO NOTHING;