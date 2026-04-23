-- name: CreateShoppingList :one
INSERT INTO shopping_lists (id, name, created_at, updated_at, user_id)
VALUES (
	gen_random_uuid(),
	$1,
	NOW(),
	NOW(),
	$2
) RETURNING *;

-- name: GetShoppingList :one
SELECT * FROM shopping_lists
WHERE id = $1;

-- name: GetUserLists :many
SELECT * FROM shopping_lists
WHERE user_id = $1;

-- name: GetListOwner :one
SELECT name, user_id FROM shopping_lists
WHERE id = $1;

-- name: PrintList :many
SELECT ingredients.name, conversions.from_unit, conversions.to_unit, conversions.ratio, shopping_list_ingredients.units, shopping_list_ingredients.quantity FROM ingredients
INNER JOIN conversions ON conversions.ingredient_id = ingredients.id
INNER JOIN shopping_list_ingredients ON shopping_list_ingredients.ingredient_id = ingredients.id
WHERE shopping_list_ingredients.shopping_list_id = $1;
