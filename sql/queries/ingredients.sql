-- name: CreateIngredient :one
INSERT INTO ingredients (id, name, image_key, created_at, updated_at)
VALUES (
	gen_random_uuid(),
	$1,
	$2,
	NOW(),
	NOW()
) RETURNING *;
