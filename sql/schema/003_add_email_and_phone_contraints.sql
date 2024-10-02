-- +goose Up

ALTER TABLE Users
ADD CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
ADD CONSTRAINT phone_number_format CHECK (phoneNumber ~* '^\+[0-9]+$');

-- +goose Down

ALTER TABLE Users
DROP CONSTRAINT email_format,
DROP CONSTRAINT phone_number_format;