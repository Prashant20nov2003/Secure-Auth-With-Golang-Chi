-- name: RegisterUser :one
INSERT INTO Users(username,email,password)
VALUES ($1,$2,$3)
RETURNING *;

-- name: LoginUser :one
SELECT * FROM Users WHERE username = $1 AND password = crypt($2, password);