-- +goose Up
CREATE TABLE shopping_lists (
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	user_id UUID NOT NULL,
	CONSTRAINT fk_users
	FOREIGN KEY (user_id)
	REFERENCES users(id)
);

-- +goose Down
DROP TABLE shopping_lists;
