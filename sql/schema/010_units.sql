-- +goose Up
CREATE TABLE units (
	name TEXT PRIMARY KEY,
	abbreviation TEXT NOT NULL
);

CREATE TABLE retail_units (
	name TEXT PRIMARY KEY
);

CREATE TABLE conversions (
	ingredient_id UUID NOT NULL,
	from_unit TEXT NOT NULL,
	to_unit TEXT NOT NULL,
	ratio REAL NOT NULL,
	FOREIGN KEY (ingredient_id) REFERENCES ingredients(id),
	FOREIGN KEY (from_unit) REFERENCES units(name),
	FOREIGN KEY (to_unit) REFERENCES retail_units(name),
	PRIMARY KEY (ingredient_id, from_unit, to_unit)
);

-- +goose Down
DROP TABLE units, retail_units, conversions;
