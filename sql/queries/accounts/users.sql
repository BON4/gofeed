-- name: CreateUser :one
INSERT INTO Users (
       uuid,
       ip,
       os,
       browser
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUser :one
SELECT * FROM Users
WHERE uuid = $1 LIMIT 1;

-- name: GetUsersByAccount :many
SELECT * FROM Users
WHERE account_username = $1;
