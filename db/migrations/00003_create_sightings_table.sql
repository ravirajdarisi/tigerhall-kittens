-- +goose Up
CREATE TABLE sightings (
  id SERIAL PRIMARY KEY,
  tiger_id INT REFERENCES tigers(id),
  lat DOUBLE PRECISION NOT NULL,
  lon DOUBLE PRECISION NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  image_path VARCHAR(255)
);

-- +goose Down
DROP TABLE sightings;
