-- name: CreateJoke :one
INSERT INTO jokes (
    author,
    title,
    text,
    explanation
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- GET QUERIES

-- name: GetJoke :one
SELECT * FROM jokes
WHERE id = $1 LIMIT 1;

-- name: ListJokes :many
SELECT * FROM jokes
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListJokesByAuthor :many
SELECT * FROM jokes
WHERE author = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- UPDATE QUERIES

-- name: UpdateJokeTitle :one
UPDATE jokes
SET title = $2
WHERE id = $1
RETURNING *;

-- name: UpdateJokeText :one
UPDATE jokes
SET text = $2
WHERE id = $1
RETURNING *;

-- name: UpdateJokeExplanation :one
UPDATE jokes
SET explanation = $2
WHERE id = $1
RETURNING *;

-- DELETE QUERIES

-- name: DeleteJoke :exec
DELETE FROM jokes
WHERE id = $1;

-- name: DeleteJokesByAuthor :exec
DELETE FROM jokes
WHERE author = $1;

-- name: DeleteAllJokes :exec
DELETE FROM jokes;