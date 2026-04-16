-- +goose Up
CREATE TABLE recipe_ingredients (
	recipe_id UUID NOT NULL,
	ingredient_id UUID NOT NULL,
	quantity REAL NOT NULL,
	unit TEXT NOT NULL,
	FOREIGN KEY (recipe_id) REFERENCES recipes(id),
	FOREIGN KEY (ingredient_id) REFERENCES ingredients(id),
	PRIMARY KEY (recipe_id, ingredient_id)
);

-- +goose Down
DROP TABLE recipe_ingredients;
