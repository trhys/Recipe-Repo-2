-- name: AddToRecipe :one
INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity, unit)
VALUES (
	$1,
	$2,
	$3,
	$4
) RETURNING *;

-- name: GetIngredientList :many
SELECT * FROM ingredients
INNER JOIN recipe_ingredients ON ingredients.id = recipe_ingredients.ingredient_id
WHERE recipe_ingredients.recipe_id = $1;
