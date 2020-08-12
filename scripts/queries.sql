-- name: GetFruit :one
SELECT * FROM fruits
WHERE id = $1;

-- name: ListFruits :many
SELECT * FROM fruits;

-- name: ListLevels :many
SELECT * FROM colorkey;

-- name: GetLevel :one
SELECT * FROM colorkey
WHERE name = '$1';

-- name: GetDetail :one
SELECT * FROM fruitinfo
WHERE name = '$1';

-- name: ListDetails :many
SELECT * FROM fruitinfo;