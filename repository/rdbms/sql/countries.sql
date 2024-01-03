CREATE TABLE IF NOT EXISTS countries
(
    country_id    SMALLINT PRIMARY KEY,
    alpha2_code   CHAR(2) UNIQUE,
    alpha3_code   CHAR(3) UNIQUE,
    olympic_code  CHAR(3) UNIQUE NULLS NOT DISTINCT,
    fifa_code     CHAR(3) UNIQUE NULLS NOT DISTINCT,
    flag          CHAR(2),
    population    INTEGER,
    area          REAL,
    independent   BOOLEAN NOT NULL,
    landlocked    BOOLEAN NOT NULL,
    un_member     BOOLEAN NOT NULL,
    latitude      REAL,
    longitude     REAL,
    region_id     INTEGER REFERENCES regions (region_id),
    subregion_id  INTEGER REFERENCES regions (region_id),
    official_name TEXT    NOT NULL,
    common_name   TEXT NOT NULL,
    start_of_week TEXT NOT NULL,
    status        TEXT NOT NULL
);
CREATE TABLE IF NOT EXISTS top_level_domains
(
    id         SERIAL PRIMARY KEY,
    country_id SMALLINT NOT NULL REFERENCES countries (country_id),
    tld        TEXT     NOT NULL,
    UNIQUE (country_id, tld)
);
CREATE TABLE IF NOT EXISTS borders
(
    id          SERIAL PRIMARY KEY,
    country_id  SMALLINT NOT NULL REFERENCES countries (country_id),
    alpha3_code CHAR(3)  NOT NULL,
    UNIQUE (country_id, alpha3_code)
);
CREATE TABLE IF NOT EXISTS country_continents
(
    country_id   SMALLINT REFERENCES countries (country_id),
    continent_id INT REFERENCES regions (region_id),
    UNIQUE (country_id, continent_id)
);

CREATE TABLE IF NOT EXISTS translations
(
    id            SERIAL PRIMARY KEY,
    country_id    SMALLINT NOT NULL REFERENCES countries (country_id),
    language_id   SMALLINT NOT NULL REFERENCES languages (language_id),
    native        bool     NOT NULL,
    official_name TEXT     NOT NULL,
    common_name   TEXT     NOT NULL,
    UNIQUE (country_id, language_id, native)
);
CREATE UNIQUE INDEX IF NOT EXISTS translations_country ON translations (country_id, language_id, native);