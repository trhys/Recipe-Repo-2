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

-- name: GetIngredientList :many
SELECT * FROM ingredients
WHERE recipe_id = $1
ORDER BY created_at;
