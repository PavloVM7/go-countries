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
    language      CHAR(3) UNIQUE     NOT NULL,
    language_name VARCHAR(32) UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS currencies
(
    currency_id   SMALLSERIAL PRIMARY KEY,
    currency      VARCHAR(3) UNIQUE  NOT NULL,
    currency_name VARCHAR(39) UNIQUE NOT NULL,
    symbol        VARCHAR(10) UNIQUE NOT NULL
);