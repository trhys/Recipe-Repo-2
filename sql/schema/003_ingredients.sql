-- +goose Up
CREATE TABLE ingredients(
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	image_key TEXT NOT NULL DEFAULT '',
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE ingredients;
