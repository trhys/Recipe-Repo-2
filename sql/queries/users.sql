-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_pw, name)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2,
	$3
) RETURNING id, created_at, email, name;

-- name: GetUserHash :one
SELECT id, hashed_pw FROM users
WHERE email = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUser :one
SELECT id, created_at, updated_at, email, name FROM USERS
WHERE id = $1;

-- name: GetName :one
SELECT name FROM users
WHERE id = $1;
