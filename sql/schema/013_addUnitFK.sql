-- +goose Up
ALTER TABLE recipe_ingredients
ADD CONSTRAINT fk_units
FOREIGN KEY (unit) REFERENCES units(name);

-- +goose Down
ALTER TABLE recipe_ingredients
DROP CONSTRAINT fk_units;
