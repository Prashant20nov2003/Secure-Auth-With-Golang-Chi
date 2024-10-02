-- name: CheckUserValidationById :one
SELECT * FROM Users WHERE user_id = $1;

-- name: CheckUserValidationByUsername :one
SELECT * FROM Users WHERE username = $1;