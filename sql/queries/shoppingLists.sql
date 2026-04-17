-- name: CreateShoppingList :one
INSERT INTO shopping_lists (id, name, created_at, updated_at, user_id)
VALUES (
	gen_random_uuid(),
	$1,
	NOW(),
	NOW(),
	$2
) RETURNING *;
