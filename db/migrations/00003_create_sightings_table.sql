-- +goose Up
CREATE TABLE sightings (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  tiger_id INT NOT NULL,
  lat DOUBLE PRECISION NOT NULL,
  lon DOUBLE PRECISION NOT NULL,
  timestamp TIMESTAMP NOT NULL,
  image_path VARCHAR(255),
  CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES users(id),
  CONSTRAINT fk_tiger FOREIGN KEY (tiger_id) REFERENCES tigers(id)
);

-- +goose Down
DROP TABLE sightings;

