-- name: CreateAccount :one
INSERT INTO Accounts (
       username,
       password,
       email,
       role
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM Accounts
WHERE username = $1 LIMIT 1;
