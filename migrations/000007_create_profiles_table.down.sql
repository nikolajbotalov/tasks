ALTER TABLE users
    ADD COLUMN username TEXT;

UPDATE users u
SET username = COALESCE(NULLIF(TRIM(p.username), ''), split_part(u.email, '@', 1))
FROM profiles p
WHERE p.user_id = u.id;

UPDATE users
SET username = split_part(email, '@', 1)
WHERE username IS NULL;

ALTER TABLE users
    ALTER COLUMN username SET NOT NULL;

DROP TABLE IF EXISTS profiles;