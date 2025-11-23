-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE teams
(
    name VARCHAR(255) NOT NULL,
    user_id VARCHAR(255) REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (name, user_id)
);

-- +goose Down
DROP TABLE IF EXISTS teams;