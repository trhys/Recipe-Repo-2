-- name: CreateRecipe :one
INSERT INTO recipes (id, created_at, updated_at, user_id)
VALUES(
	gen_random_uuid(),
	NOW(),
	NOW(),
	$1
)
RETURNING *;
