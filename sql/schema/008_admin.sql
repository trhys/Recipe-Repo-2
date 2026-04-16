-- +goose Up
ALTER TABLE users
ADD COLUMN admin BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE users
DROP COLUMN admin;
