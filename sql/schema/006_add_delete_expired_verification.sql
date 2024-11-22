-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION delete_expired_verification()
RETURNS void
LANGUAGE plpgsql
AS $$
BEGIN
  DELETE FROM UserEmailVerifications
  WHERE expires_at < CURRENT_TIMESTAMP;
END;
$$;
-- +goose StatementEnd

-- +goose Down
DROP FUNCTION IF EXISTS delete_expired_verification();