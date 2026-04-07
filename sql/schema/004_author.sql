-- +goose Up
ALTER TABLE users
ADD COLUMN name TEXT NOT NULL DEFAULT 'anonymous';

-- +goose Down
ALTER TABLE users
DROP COLUMN name;
