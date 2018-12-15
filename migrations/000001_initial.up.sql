BEGIN;

CREATE TABLE users (
    id      SERIAL PRIMARY KEY,
    gid     NUMERIC NOT NULL,
    name    TEXT    NOT NULL,
    email   TEXT    NOT NULL,
    picture TEXT
);

CREATE INDEX ON users (gid);

COMMIT;
