CREATE TABLE IF NOT EXISTS regions
(
    region_id   SERIAL PRIMARY KEY,
    parent_id   INTEGER NOT NULL,
    region_name TEXT    NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS regions_parent ON regions (parent_id, region_name);

CREATE TABLE IF NOT EXISTS languages
(
    language_id SMALLSERIAL PRIMARY KEY,
    language    CHAR(3) UNIQUE NOT NULL
);