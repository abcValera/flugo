-- name: CreateUser :one
INSERT INTO users(
    username,
    email,
    hashed_password,
    fullname,
    status,
    bio
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- GET QUERIES

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByName :one
SELECT * FROM users
WHERE username = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ListUsers :many
SELECT * FROM users
ORDER BY id
LIMIT $1
OFFSET $2;

-- UPDATE QUERIES

-- name: UpdateUserPassword :one
UPDATE users
SET hashed_password = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserFullname :one
UPDATE users
SET fullname = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserStatus :one
UPDATE users
SET status = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserBio :one
UPDATE users
SET bio = $2
WHERE id = $1
RETURNING *;

-- name: UpdateUserUpdatedAt :exec
UPDATE users
SET updated_at = now()
WHERE id = $1;

-- DELETE QUERIES

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;
