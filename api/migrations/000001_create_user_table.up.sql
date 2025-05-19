CREATE TABLE users
(
    username VARCHAR(255) NOT NULL,
    login    VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL
);
