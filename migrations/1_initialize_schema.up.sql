BEGIN;


CREATE TABLE IF NOT EXISTS bookmarks (
    id          SERIAL PRIMARY KEY,
    title       varchar(100),
    description TEXT,
    body        TEXT,
    image       varchar(100),
    url         varchar(150),
    archived    boolean,
    deleted     boolean,
    created_at  timestamp(6) with time zone DEFAULT CURRENT_TIMESTAMP(6),
    updated_at  timestamp(6) with time zone DEFAULT NULL::timestamp(6) with time zone
);

COMMIT;