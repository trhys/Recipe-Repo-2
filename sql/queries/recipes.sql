-- name: CreateRecipe :one
INSERT INTO recipes (id, title, created_at, updated_at, user_id)
VALUES(
	gen_random_uuid(),
	$1,
	NOW(),
	NOW(),
	$2
)
RETURNING *;

-- name: GetRecipe :one
SELECT * FROM recipes
WHERE id = $1;
