CREATE TABLE role_permissions
(
    role       TEXT NOT NULL,
    permission TEXT NOT NULL,
    PRIMARY KEY (role, permission),
    FOREIGN KEY (role) REFERENCES roles (name) ON DELETE CASCADE ON UPDATE CASCADE
);