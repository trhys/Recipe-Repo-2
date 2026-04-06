-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_pw)
VALUES (
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1,
	$2
) RETURNING id, created_at, email;

-- name: GetUserHash :one
SELECT id, hashed_pw FROM users
WHERE email = $1;

-- name: ResetUsers :exec
DELETE FROM users;
