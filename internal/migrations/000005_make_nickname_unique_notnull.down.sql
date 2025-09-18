ALTER TABLE users
DROP CONSTRAINT users_nickname_unique;

ALTER TABLE users
ALTER COLUMN nickname DROP NOT NULL;
