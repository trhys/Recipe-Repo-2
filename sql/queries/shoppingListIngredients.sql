-- name: AddToShoppingList :exec
INSERT INTO shopping_list_ingredients (shopping_list_id, ingredient_id, quantity, units)
VALUES (
	$1,
	$2,
	$3,
	$4
);
