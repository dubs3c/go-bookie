BEGIN;

CREATE TABLE IF NOT EXISTS bookmarks (
    id          SERIAL PRIMARY KEY,
    title       TEXT DEFAULT '',
    description TEXT DEFAULT '',
    body        TEXT DEFAULT '',
    image       VARCHAR(150) DEFAULT '',
    url         TEXT NOT NULL,
    archived    BOOLEAN DEFAULT FALSE,
    deleted     BOOLEAN DEFAULT FALSE,
    created_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  timestamp(6) with time zone DEFAULT NULL::timestamp(6) with time zone
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