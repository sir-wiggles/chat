BEGIN;
ALTER TABLE users RENAME username TO name;
COMMIT;
