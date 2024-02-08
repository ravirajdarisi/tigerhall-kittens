-- +goose Up
CREATE TABLE tigers (
  id SERIAL PRIMARY KEY,
  name VARCHAR(255) NOT NULL,
  date_of_birth DATE,
  last_seen_timestamp TIMESTAMP,
  last_seen_lat DOUBLE PRECISION,
  last_seen_lon DOUBLE PRECISION
);

-- +goose Down
DROP TABLE tigers;
