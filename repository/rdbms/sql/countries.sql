CREATE TABLE IF NOT EXISTS countries
(
    country_id    SMALLINT PRIMARY KEY,
    alpha2_code   CHAR(2) UNIQUE,
    alpha3_code   CHAR(3) UNIQUE,
    olympic_code  CHAR(3) UNIQUE NULLS NOT DISTINCT,
    fifa_code     CHAR(3) UNIQUE NULLS NOT DISTINCT,
    flag          CHAR(4),
    population    INTEGER,
    area          REAL,
    independent   BOOLEAN NOT NULL,
    landlocked    BOOLEAN NOT NULL,
    unMember      BOOLEAN NOT NULL,
    latitude      REAL,
    longitude     REAL,
    region_id     INTEGER REFERENCES regions (region_id),
    subregion_id  INTEGER REFERENCES regions (region_id),
    official_name TEXT,
    common_name   TEXT
);

CREATE TABLE IF NOT EXISTS country_continents
(
    country_id   SMALLINT REFERENCES countries (country_id),
    continent_id INT REFERENCES regions (region_id),
    UNIQUE (country_id, country_id)
);