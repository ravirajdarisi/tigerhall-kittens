-- +goose Up
CREATE TABLE tigers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL UNIQUE, -- This enforces uniqueness for the 'name' column
  date_of_birth DATE NOT NULL,
  last_seen_timestamp TIMESTAMP NOT NULL,
  last_seen_lat DOUBLE PRECISION NOT NULL,
  last_seen_lon DOUBLE PRECISION NOT NULL
);


-- +goose Down
DROP TABLE tigers;
