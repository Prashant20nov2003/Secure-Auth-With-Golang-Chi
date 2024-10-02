-- +goose Up
CREATE TABLE Products(
  product_id VARCHAR(40) PRIMARY KEY DEFAULT CONCAT('PR-', uuid_generate_v4()),
  user_id VARCHAR(40) NOT NULL REFERENCES Users(user_id) ON DELETE CASCADE ON UPDATE CASCADE,
  product_name VARCHAR(100) NOT NULL,
  product_photo UUID NOT NULL DEFAULT uuid_generate_v4() UNIQUE,
  price INT NOT NULL,
  visibility BOOLEAN NOT NULL DEFAULT true
);

-- +goose Down
DROP TABLE Products;