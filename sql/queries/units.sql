-- name: CreateConversion :exec
INSERT INTO conversions (ingredient_id, from_unit, to_unit, ratio)
VALUES (
	$1,
	$2,
	$3,
	$4
);
