-- +goose Up
CREATE TABLE shopping_list_recipes (
	shopping_list_id UUID NOT NULL,
	recipe_id UUID NOT NULL,
	quantity INT NOT NULL,
	FOREIGN KEY (shopping_list_id) REFERENCES shopping_lists(id) ON DELETE CASCADE,
	FOREIGN KEY (recipe_id) REFERENCES recipes(id),
	PRIMARY KEY (shopping_list_id, recipe_id)
);

-- +goose Down
DROP TABLE shopping_list_recipes;
