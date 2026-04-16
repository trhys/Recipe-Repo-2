-- +goose Up
CREATE TABLE users(
	id UUID PRIMARY KEY,
	name TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,
	updated_at TIMESTAMP NOT NULL,
	email TEXT NOT NULL UNIQUE,
	hashed_pw TEXT NOT NULL
);

-- +goose Down
DROP TABLE users;
