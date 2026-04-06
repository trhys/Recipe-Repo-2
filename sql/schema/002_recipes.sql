-- +goose Up
CREATE TABLE recipes(
	id UUID PRIMARY KEY,
	title TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID NOT NULL,
	FOREIGN KEY (user_id)
	REFERENCES users(id)
	ON DELETE CASCADE
);

-- +goose Down
DROP TABLE recipes;
