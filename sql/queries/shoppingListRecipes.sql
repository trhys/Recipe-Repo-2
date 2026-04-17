-- name: AddRecipeToList :exec
INSERT INTO shopping_list_recipes (shopping_list_id, recipe_id, quantity)
VALUES (
	$1,
	$2,
	$3
);

-- name: GetRecipesFromList :many
SELECT recipes.*, shopping_list_recipes.quantity FROM recipes
INNER JOIN shopping_list_recipes ON shopping_list_recipes.recipe_id = recipes.id
WHERE shopping_list_recipes.shopping_list_id = $1;
