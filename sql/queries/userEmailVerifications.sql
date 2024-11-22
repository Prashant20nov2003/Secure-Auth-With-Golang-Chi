-- name: CreateEmailVerification :one
SELECT
user_id::VARCHAR,
emailverify_id::UUID
FROM generate_email_verification($1,$2,$3)
AS result(user_id,emailverify_id);

-- name: FetchEmailVerification :one
SELECT
  EXTRACT(EPOCH FROM (expires_at - CURRENT_TIMESTAMP))::INTEGER as time_left
FROM useremailverifications
WHERE emailverify_id = $1;

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

-- name: ResendEmailVerificationCode :one
UPDATE UserEmailVerifications
SET verif_code = $1, expires_at = CURRENT_TIMESTAMP + INTERVAL '60 seconds'
FROM users
WHERE UserEmailVerifications.emailverify_id = $2
  AND users.user_id = UserEmailVerifications.user_id
RETURNING
UserEmailVerifications.emailverify_id,
EXTRACT(EPOCH FROM (UserEmailVerifications.expires_at - CURRENT_TIMESTAMP))::INTEGER as time_left,
users.email;