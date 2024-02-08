package models

import (
	"database/sql"
	_"errors"
	"time"
)

// Sighting struct represents an observation of a tiger in the wild
type Sighting struct {
	ID        int       `json:"id"`
	TigerID   int       `json:"tiger_id"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	Timestamp time.Time `json:"timestamp"`
	ImagePath string    `json:"image_path"`
}

// NewSighting creates a new Sighting instance.
func NewSighting(tigerID int, lat, lon float64, imagePath string) *Sighting {
	return &Sighting{
		TigerID:   tigerID,
		Lat:       lat,
		Lon:       lon,
		Timestamp: time.Now(),
		ImagePath: imagePath,
	}
}

// Save inserts the Sighting into the database.
func (s *Sighting) Save(db *sql.DB) error {
	query := `INSERT INTO sightings (tiger_id, lat, lon, timestamp, image_path) VALUES ($1, $2, $3, $4, $5)`
	_, err := db.Exec(query, s.TigerID, s.Lat, s.Lon, s.Timestamp, s.ImagePath)
	return err
}

// GetLastSightingByTigerID retrieves the most recent sighting of a given tiger from the database.
// Todo : should work on the logic to check 5km radius
func GetLastSightingByTigerID(db *sql.DB, tigerID int) (*Sighting, error) {
	query := `SELECT id, tiger_id, lat, lon, timestamp, image_path FROM sightings WHERE tiger_id = $1 ORDER BY timestamp DESC LIMIT 1`
	sighting := Sighting{}
	err := db.QueryRow(query, tigerID).Scan(&sighting.ID, &sighting.TigerID, &sighting.Lat, &sighting.Lon, &sighting.Timestamp, &sighting.ImagePath)
	if err != nil {
		if err == sql.ErrNoRows {
			// No sightings is not an error , so should return nil 
			return nil, nil
		}
		return nil, err
	}
	return &sighting, nil
}

// GetAllSightingsByTigerID retrieves all sightings of a given tiger from the database.
// useful method to have a history
func GetAllSightingsByTigerID(db *sql.DB, tigerID int) ([]Sighting, error) {
	query := `SELECT id, tiger_id, lat, lon, timestamp, image_path FROM sightings WHERE tiger_id = $1 ORDER BY timestamp DESC`
	rows, err := db.Query(query, tigerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sightings []Sighting
	for rows.Next() {
		var sighting Sighting
		if err := rows.Scan(&sighting.ID, &sighting.TigerID, &sighting.Lat, &sighting.Lon, &sighting.Timestamp, &sighting.ImagePath); err != nil {
			return nil, err
		}
		sightings = append(sightings, sighting)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return sightings, nil
}
