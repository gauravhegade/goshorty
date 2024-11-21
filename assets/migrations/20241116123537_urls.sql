-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS urls (
    short_code TEXT PRIMARY KEY,
    long_url TEXT NOT NULL,
    title TEXT,
    created_on DATETIME NOT NULL DEFAULT (DATETIME('now')),
    expires_on DATETIME
);

INSERT INTO urls (short_code, long_url, title, created_on, expires_on) 
VALUES 
('abc123', 'https://example.com/long-url-example-1', 'Example 1', DATETIME('now'), DATETIME('now', '+30 days'));

INSERT INTO urls (short_code, long_url, title, created_on, expires_on) 
VALUES 
('xyz789', 'https://example.com/long-url-example-2', 'Example 2', DATETIME('now'), DATETIME('now', '+60 days'));

INSERT INTO urls (short_code, long_url, title, created_on, expires_on) 
VALUES 
('pqr456', 'https://example.com/long-url-example-3', 'Example 3', DATETIME('now'), DATETIME('now', '+90 days'));

INSERT INTO urls (short_code, long_url, title, created_on, expires_on) 
VALUES 
('lmn321', 'https://example.com/long-url-example-4', 'Example 4', DATETIME('now'), DATETIME('now', '+7 days'));

INSERT INTO urls (short_code, long_url, title, created_on, expires_on) 
VALUES 
('stu654', 'https://example.com/long-url-example-5', 'Example 5', DATETIME('now'), NULL);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS ShortURLAliases;
DROP INDEX IF EXISTS ShortURLKeys;
DROP TABLE IF EXISTS urls;
-- +goose StatementEnd