-- name: ChangePassword :one
SELECT 
  user_id::VARCHAR,
  username::VARCHAR,
  email::VARCHAR,
  isEmailVerified::VARCHAR
FROM change_password($1, $2)
AS result(user_id, username, email, isEmailVerified);