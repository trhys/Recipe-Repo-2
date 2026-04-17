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
