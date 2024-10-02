-- +goose Up

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE Users(
  user_id VARCHAR(40) PRIMARY KEY DEFAULT CONCAT('ID-', uuid_generate_v4()),
  username VARCHAR(20) NOT NULL UNIQUE,
  email VARCHAR(50) NOT NULL UNIQUE,
  phoneNumber VARCHAR(20) NOT NULL UNIQUE,
  password VARCHAR(60) NOT NULL,
  last_login TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE Users;

DROP EXTENSION IF EXISTS "uuid-ossp";