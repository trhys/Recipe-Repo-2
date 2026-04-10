-- +goose Up
ALTER TABLE recipes
ADD COLUMN image_key TEXT NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE recipes
DROP COLUMN image_key;
