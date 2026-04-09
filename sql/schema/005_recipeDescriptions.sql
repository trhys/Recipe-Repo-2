-- +goose Up
ALTER TABLE recipes
ADD COLUMN description TEXT NOT NULL DEFAULT 'No description provided...';

-- +goose Down
ALTER TABLE recipes
DROP COLUMN description;
