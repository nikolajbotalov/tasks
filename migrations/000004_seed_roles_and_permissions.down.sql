DELETE FROM role_permissions WHERE role IN ('admin', 'user', 'viewer');
DELETE FROM roles WHERE name IN ('admin', 'user', 'viewer');