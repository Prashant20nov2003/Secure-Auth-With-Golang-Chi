-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION generate_email_verification(
  _user_id VARCHAR, _verif_code VARCHAR, _used_for VARCHAR
)
RETURNS TABLE (
  user_id VARCHAR,
  emailverify_id UUID
)
LANGUAGE plpgsql
AS $$
BEGIN
  RETURN QUERY
    INSERT INTO useremailverifications(
      user_id, verif_code, used_for)
    VALUES (_user_id, _verif_code, _used_for)
  RETURNING useremailverifications.user_id, useremailverifications.emailverify_id;
END;
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS generate_email_verification;
