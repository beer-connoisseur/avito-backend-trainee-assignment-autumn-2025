-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE pull_requests
(
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL,
    author_id VARCHAR(255) REFERENCES users (id),
    status VARCHAR(20) CHECK (status IN ('OPEN', 'MERGED')) DEFAULT 'OPEN',
    assigned_reviewer_first VARCHAR(255) REFERENCES users (id),
    assigned_reviewer_second VARCHAR(255) REFERENCES users (id),
    created_at TIMESTAMP DEFAULT NOW(),
    merged_at TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS pull_requests;