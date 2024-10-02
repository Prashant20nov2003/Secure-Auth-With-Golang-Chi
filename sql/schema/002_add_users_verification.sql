-- +goose Up

ALTER TABLE Users
ALTER COLUMN last_login SET DEFAULT NOW(),
ALTER COLUMN phoneNumber SET DEFAULT NULL,
ALTER COLUMN phoneNumber DROP NOT NULL;


-- +goose Down

ALTER TABLE Users
ALTER COLUMN last_login DROP DEFAULT,
ALTER COLUMN phoneNumber DROP DEFAULT,
ALTER COLUMN phoneNumber SET NOT NULL;
