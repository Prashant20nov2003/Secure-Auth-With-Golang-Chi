-- name: PostProduct :one
INSERT INTO products(user_id, product_name, price, visibility)
VALUES ($1,$2,$3,$4)
RETURNING *;

-- name: GetPublicProduct :many
SELECT *, v.username FROM products u
JOIN users v ON v.user_id = u.user_id
WHERE u.user_id!=$1 AND u.visibility=true;

-- name: GetMyProduct :many
SELECT * FROM products WHERE user_id=$1;

-- name: UpdateProduct :one
-- UPDATE Products
-- SET
-- product_name = COALESCE($1,product_name),
-- product_photo = COALESCE($2,product_photo),
-- price = COALESCE($3,product_photo),
-- visibility = COALESCE($4,visibility)
-- WHERE product_id=$5 AND user_id=$5
-- RETURNING *;