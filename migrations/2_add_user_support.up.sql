BEGIN;

-- For generating cryptographically secure bytes with gen_random_bytes
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,
    email       VARCHAR(100) NOT NULL UNIQUE,
    password    VARCHAR(100) NOT NULL,
    api_token   VARCHAR(64) DEFAULT encode(sha256(gen_random_bytes(32)), 'hex'),
    api_enabled boolean DEFAULT false,
    is_admin    boolean DEFAULT false,
    created_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6)
);

CREATE TABLE IF NOT EXISTS access_tokens (
    id          SERIAL PRIMARY KEY,
    token       VARCHAR(64) UNIQUE DEFAULT encode(sha256(gen_random_bytes(32)), 'hex'),
    user_fk     INTEGER,
    created_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE
);

-- useful if applying to DB where data already exists
INSERT INTO users(email, password) VALUES('admin@samla.app', '1234');
ALTER TABLE IF EXISTS bookmarks ADD COLUMN user_fk INTEGER NOT NULL DEFAULT 1;


ALTER TABLE IF EXISTS bookmarks ADD CONSTRAINT fk_users_bookmarks FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE;
ALTER TABLE IF EXISTS bookmarks ALTER COLUMN user_fk DROP DEFAULT;

/*CREATE TABLE IF NOT EXISTS user_has_bookmarks (
    id              SERIAL PRIMARY KEY,
    user_fk         INTEGER,
    bookmark_fk     INTEGER,
    FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (bookmark_fk) REFERENCES bookmarks(id) ON DELETE CASCADE,
    UNIQUE(user_fk, bookmark_fk)
);*/

COMMIT;