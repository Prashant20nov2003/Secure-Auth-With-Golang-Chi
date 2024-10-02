-- +goose Up
ALTER TABLE UserEmailVerifications
ADD COLUMN used_for VARCHAR(30) NOT NULL,
ADD COLUMN is_verified BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE UserEmailVerifications
DROP COLUMN used_for,
DROP COLUMN is_verified;