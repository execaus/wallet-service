-- name: Get :one
SELECT *
FROM app.wallets
WHERE id = $1;

-- name: GetForUpdate :one
SELECT *
FROM app.wallets
WHERE id = $1
FOR UPDATE;

-- name: Update :one
UPDATE app.wallets
SET balance = $2
WHERE id = $1
RETURNING *;