-- name: CreateEmailVerification :one
SELECT
user_id::VARCHAR,
emailverify_id::UUID
FROM generate_email_verification($1,$2,$3)
AS result(user_id,emailverify_id);

-- name: CheckEmailVerification :one
SELECT 
  res_message::TEXT,
  emailverify_id::UUID,
  user_id::VARCHAR,
  username::VARCHAR,
  email::VARCHAR,
  used_for::VARCHAR
FROM check_email_verification($1, $2)
AS result(res_message, emailverify_id, user_id, username, email, used_for);

-- name: UpdateUserEmail :one
SELECT
  message::VARCHAR,
  user_id::VARCHAR,
  username::VARCHAR,
  email::VARCHAR
FROM update_user_email($1, $2)
AS result(message,user_id,username,email);