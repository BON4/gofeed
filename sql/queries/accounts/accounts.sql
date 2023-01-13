-- name: CreateAccount :one
INSERT INTO Accounts (
       username,
       password,
       email,
       role,
       activated
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM Accounts
WHERE username = $1 LIMIT 1;
