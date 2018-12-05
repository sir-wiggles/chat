BEGIN;
ALTER TABLE users RENAME name TO username;
COMMIT;
