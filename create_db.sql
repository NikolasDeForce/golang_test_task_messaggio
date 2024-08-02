DROP DATABASE IF EXISTS messages_rest;
CREATE DATABASE messages_rest;

\c messages_rest;

CREATE TABLE lists (
id SERIAL PRIMARY KEY,
text VARCHAR NOT NULL
);

INSERT INTO lists(text) VALUES ('hello world!');