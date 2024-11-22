-- +goose Up
CREATE TABLE UserEmailVerifications(
  emailverify_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  user_id VARCHAR(40) NOT NULL REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
  expires_at TIMESTAMP NOT NULL DEFAULT (CURRENT_TIMESTAMP + INTERVAL '60 seconds'),
  verif_code VARCHAR(60) NOT NULL,
  attempts SMALLINT NOT NULL DEFAULT 0
);

-- +goose Down
DROP TABLE UserEmailVerifications;