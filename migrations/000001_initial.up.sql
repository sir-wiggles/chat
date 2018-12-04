BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id     SERIAL PRIMARY KEY,
    uuid   UUID NOT NULL DEFAULT uuid_generate_v4(),
    name   TEXT NOT NULL,
    email  TEXT,
    avatar TEXT
);

CREATE INDEX ON users (name);

COMMIT;
