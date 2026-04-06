-- name: CreateIngredient :one
INSERT INTO ingredients (id, name, quantity, unit, created_at, updated_at, recipe_id)
VALUES (
	gen_random_uuid(),
	$1,
	$2,
	$3,
	NOW(),
	NOW(),
	$4
) RETURNING *;
