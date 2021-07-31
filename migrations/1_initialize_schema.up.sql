BEGIN;


CREATE TABLE IF NOT EXISTS bookmarks (
    id          SERIAL PRIMARY KEY,
    title       VARCHAR(100) DEFAULT "",
    description TEXT DEFAULT "",
    body        TEXT DEFAULT "",
    image       VARCHAR(100) DEFAULT "",
    url         VARCHAR(150) NOT NULL,
    archived    BOOLEAN DEFAULT FALSE,
    deleted     BOOLEAN DEFAULT FALSE,
    created_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  timestamp(6) with time zone DEFAULT NULL::timestamp(6) with time zone
);

COMMIT;