-- +goose Up
CREATE TABLE shopping_list_ingredients (
	shopping_list_id UUID NOT NULL,
	ingredient_id UUID NOT NULL,
	quantity REAL NOT NULL,
	FOREIGN KEY (shopping_list_id) REFERENCES shopping_lists(id),
	FOREIGN KEY (ingredient_id) REFERENCES ingredients(id),
	PRIMARY KEY (shopping_list_id, ingredient_id)
);

-- +goose Down
DROP TABLE shopping_list_ingredients;
