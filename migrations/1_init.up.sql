CREATE TABLE IF NOT EXISTS users
(
    id           pg_catalog.uuid PRIMARY KEY,
    email        TEXT NOT NULL UNIQUE,
    fullName     TEXT NOT NULL,
    passwordHash TEXT NOT NULL,
    role         TEXT not null default 'user'
);
CREATE INDEX IF NOT EXISTS idx_email ON users (email);