-- +goose Up
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION check_email_verification(
  _verif_code VARCHAR,
  _emailverify_id UUID
)
RETURNS TABLE (
  res_message TEXT,
  emailverify_id UUID,
  user_id VARCHAR,
  username VARCHAR,
  email VARCHAR,
  used_for VARCHAR
)
LANGUAGE plpgsql
AS $$
DECLARE
  email_verification RECORD;
BEGIN
  -- Select matched useremailverification
  WITH emailverification AS (
    SELECT
      CASE
        WHEN u.expires_at < CURRENT_TIMESTAMP THEN
          'Code has been expired, please try again'
        WHEN u.attempts >= 3 THEN
          'Maximum attempts reached'
        WHEN u.verif_code != crypt(_verif_code, u.verif_code) THEN
          'Wrong verification code, please try again (' || (3 - u.attempts - 1) || ' time left)'
        ELSE
          'Success'
      END AS res_message,
      u.emailverify_id,
      u.user_id,
      v.username,
      v.email,
      u.used_for
    FROM useremailverifications u
    JOIN users v ON v.user_id = u.user_id
    WHERE u.emailverify_id = _emailverify_id
  )
  SELECT * INTO email_verification FROM emailverification;

  -- Check if we got exactly one row
  IF email_verification IS NOT NULL THEN
    IF email_verification.res_message = 'Success' THEN
      UPDATE UserEmailVerifications u
      SET is_verified = true
      WHERE u.emailverify_id = _emailverify_id;
      RETURN QUERY SELECT email_verification.res_message,
      email_verification.emailverify_id,
      email_verification.user_id,
      email_verification.username,
      email_verification.email,
      email_verification.used_for;
    ELSE
      -- Update attempts if the message indicates a wrong verification code
      IF email_verification.res_message LIKE 'Wrong verification code%' THEN
        UPDATE useremailverifications
        SET attempts = attempts + 1
        WHERE email_verification.emailverify_id = _emailverify_id;
      END IF;
      RETURN QUERY SELECT email_verification.res_message, NULL::UUID, 'NULL'::VARCHAR, 'NULL'::VARCHAR, 'NULL'::VARCHAR, 'NULL'::VARCHAR;
    END IF;
  ELSE
    -- Handle the case where no row is found
    RETURN QUERY SELECT 'No verification entry found'::TEXT, NULL::UUID, 'NULL'::VARCHAR, 'NULL'::VARCHAR, 'NULL'::VARCHAR, 'NULL'::VARCHAR;
  END IF;
END;
$$
-- +goose StatementEnd

-- +goose Down
DROP EXTENSION IF EXISTS pgcrypto;

DROP FUNCTION check_email_verification;