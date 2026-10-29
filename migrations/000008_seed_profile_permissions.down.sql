DELETE
FROM role_permissions
WHERE (role, permission) IN (
                             ('admin', 'profile:read'),
                             ('admin', 'profile:update'),
                             ('user', 'profile:read'),
                             ('user', 'profile:update'),
                             ('viewer', 'profile:read'),
                             ('viewer', 'profile:update')
    );