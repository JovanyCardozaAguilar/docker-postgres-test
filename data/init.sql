CREATE TABLE Account (
    id SERIAL PRIMARY KEY,
    FirstName VARCHAR(255) NOT NULL,
    LastName VARCHAR(255) NOT NULL,
    Token VARCHAR(255) NOT NULL
);

CREATE TABLE temp (
    id SERIAL PRIMARY KEY,
    FirstName VARCHAR(255) NOT NULL,
    LastName VARCHAR(255) NOT NULL,
    Token VARCHAR(255) NOT NULL
);

COPY temp
FROM '/docker-entrypoint-initdb.d/test.csv'
DELIMITER ','
CSV HEADER;

INSERT INTO Account (id, FirstName, LastName, Token)
SELECT id, FirstName, LastName, Token
FROM temp;

DROP TABLE IF EXISTS temp;
