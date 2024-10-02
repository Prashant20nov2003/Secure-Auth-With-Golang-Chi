-- +goose Up
ALTER TABLE Users
ADD COLUMN isEmailVerified BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN isPhoneVerified BOOLEAN NOT NULL DEFAULT false;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_user_email(
  _user_id VARCHAR,
  _emailverify_id UUID
)
RETURNS TABLE (
  message VARCHAR,
  user_id VARCHAR,
  username VARCHAR,
  email VARCHAR
)
LANGUAGE plpgsql
AS $$
  DECLARE
    user_record RECORD;
  BEGIN
    -- Update User Email
    UPDATE Users
    SET isEmailVerified = true
    WHERE Users.user_id = _user_id
    RETURNING Users.user_id, Users.username, Users.email INTO user_record;

    -- Delete The Verification Record
    DELETE FROM useremailverifications
    WHERE useremailverifications.emailverify_id = _emailverify_id;

    RETURN QUERY SELECT 'Successful verify user email'::VARCHAR,
      user_record.user_id,
      user_record.username,
      user_record.email;
    
  END;
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS update_user_email(
  _user_id VARCHAR,
  _emailverify_id UUID
);

ALTER TABLE Users
DROP COLUMN isEmailVerified,
DROP COLUMN isPhoneVerified;