-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION change_password(
  _password VARCHAR,
  _emailverify_id UUID
)
RETURNS TABLE(
  user_id VARCHAR,
  username VARCHAR,
  email VARCHAR,
  isEmailVerified BOOLEAN
)
LANGUAGE plpgsql
AS $$
DECLARE
  v_user_id VARCHAR;
BEGIN
  DELETE FROM UserEmailVerifications v
  WHERE v.emailverify_id = _emailverify_id AND v.used_for = 'Change Password' AND v.is_verified = true
  RETURNING v.user_id INTO v_user_id;

  IF v_user_id IS NULL THEN
    RETURN QUERY SELECT 'NULL'::VARCHAR,'NULL'::VARCHAR,'NULL'::VARCHAR,false::BOOLEAN;
  ELSE
    RETURN QUERY
      UPDATE Users u
      SET password = _password
      WHERE u.user_id = v_user_id
      RETURNING u.user_id, u.username, u.email, u.isEmailVerified;
  END IF;
END;
$$
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION change_password;