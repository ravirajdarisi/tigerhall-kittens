package models

import (
	"database/sql"
	_ "errors"
	"time"
)

// Sighting struct represents an observation of a tiger in the wild
type Sighting struct {
    ID         int       	`json:"id"` 
    UserID     int      	`json:"user_id"`
	TigerID    int          `json:"tiger_id"`
    Lat        float64   	`json:"lat"` 
    Lon        float64   	`json:"lon"`
    Timestamp  time.Time 	`json:"timestamp"`
    ImagePath  string   	`json:"image_path"` 
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
	query := `INSERT INTO sightings (user_id,tiger_id, lat, lon, timestamp, image_path) VALUES ($1, $2, $3, $4, $5,$6)`
	_, err := db.Exec(query, s.TigerID, s.Lat, s.Lon, s.Timestamp, s.ImagePath)
	return err
}



// GetLastSightingByTigerID retrieves the most recent sighting of a given tiger from the database.
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


// GetUsersByTigerID retrieves all unique user IDs who have reported a sighting of the given tiger.
func GetUsersByTigerID(db *sql.DB, tigerID int) ([]int, error) {
    // Prepare the SQL query to select distinct user IDs where the tiger_id matches.
    query := `SELECT DISTINCT user_id FROM sightings WHERE tiger_id = $1`
    
    // Execute the query.
    rows, err := db.Query(query, tigerID)
    if err != nil {
        // Handle any errors that occur during query execution.
        return nil, err
    }
    defer rows.Close() // Ensure rows are closed after the function finishes.

    var userIDs []int // Slice to hold the user IDs.
    for rows.Next() {
        var userID int
        if err := rows.Scan(&userID); err != nil {
            // Handle any errors that occur during row scanning.
            return nil, err
        }
        userIDs = append(userIDs, userID) // Append the user ID to the slice.
    }

    // Check for errors from iterating over rows.
    if err = rows.Err(); err != nil {
        // Handle any errors that occurred during the iteration.
        return nil, err
    }

    return userIDs, nil 
}



// GetAllSightingsByTigerID retrieves all sightings of a given tiger from the database with pagination.
func GetAllSightingsByTigerID(db *sql.DB, tigerID, limit, offset int) ([]Sighting, error) {
	query := `SELECT id, tiger_id, lat, lon, timestamp, image_path FROM sightings WHERE tiger_id = $1 ORDER BY timestamp DESC LIMIT $2 OFFSET $3`
	rows, err := db.Query(query, tigerID, limit, offset)
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
