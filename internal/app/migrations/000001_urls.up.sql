CREATE TABLE IF NOT EXISTS urls
(
    user_id      VARCHAR(36),
    short_url    VARCHAR(255) PRIMARY KEY,
    original_url VARCHAR(255),
    deleted      BOOLEAN DEFAULT false
);