-- +goose Up
CREATE TABLE ingredients(
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	quantity REAL NOT NULL,
	unit TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	recipe_id UUID NOT NULL,
	FOREIGN KEY (recipe_id)
	REFERENCES recipes(id)
	ON DELETE CASCADE
);

-- +goose Down
DROP TABLE ingredients;
