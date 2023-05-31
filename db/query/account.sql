-- name: CreateAccount :exec
INSERT INTO accounts (owner,
                      balance,
                      currency)
VALUES ($1, $2, $3) RETURNING *;