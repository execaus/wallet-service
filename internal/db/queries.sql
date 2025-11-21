-- name: Get :one
SELECT *
FROM wallets
WHERE id = $1;

-- name: GetForUpdate :one
SELECT *
FROM wallets
WHERE id = $1
FOR UPDATE;

-- name: Update :exec
UPDATE wallets
SET balance = $2
WHERE id = $1;