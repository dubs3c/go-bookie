BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id              SERIAL PRIMARY KEY,
    email           VARCHAR(50) UNIQUE NOT NULL,
    password        VARCHAR(72) NOT NULL,
    name            VARCHAR(60) DEFAULT '',
    is_superuser    BOOLEAN,
    is_active       BOOLEAN,
    created_at      timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    last_login      timestamp(6) with time zone DEFAULT NULL::timestamp(6) with time zone
);

CREATE TABLE IF NOT EXISTS access_tokens (
    id          SERIAL PRIMARY KEY,
    token       VARCHAR(100) UNIQUE NOT NULL,
    user_fk     INTEGER,
    created_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bookmarks (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(100) DEFAULT "",
    description TEXT DEFAULT "",
    body        TEXT DEFAULT "",
    image       VARCHAR(100) DEFAULT "",
    url         VARCHAR(150) NOT NULL,
    archived    BOOLEAN DEFAULT FALSE,
    deleted     BOOLEAN DEFAULT FALSE,
    user_fk     INTEGER NOT NULL,
    created_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  timestamp(6) with time zone DEFAULT NULL::timestamp(6) with time zone,
    FOREIGN KEY (user_fk) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tags (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) NOT NULL UNIQUE,
    created_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  timestamp(6) with time zone DEFAULT NULL::timestamp(6) with time zone
);

CREATE TABLE IF NOT EXISTS bookmark_has_tags (
    id          SERIAL PRIMARY KEY,
    bookmark_fk INTEGER,
    tag_fk      INTEGER,
    FOREIGN KEY (bookmark_fk) REFERENCES bookmarks(id) ON DELETE CASCADE,
    FOREIGN KEY (tag_fk) REFERENCES tags(id) ON DELETE CASCADE,
    UNIQUE(bookmark_fk, tag_fk)
);

COMMIT;