CREATE TABLE IF NOT EXISTS regions
(
    region_id   SERIAL PRIMARY KEY,
    region_name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS languages
(
    language_id SMALLSERIAL PRIMARY KEY,
    language    CHAR(3) UNIQUE NOT NULL
);