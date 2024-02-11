package models

import (
	"database/sql"
	"time"
)

// Tiger represents the tiger structure.
type Tiger struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	DateOfBirth       time.Time `json:"date_of_birth"`
	LastSeenTimestamp time.Time `json:"last_seen_timestamp"`
	LastSeenLat       float64   `json:"last_seen_lat"`
	LastSeenLon       float64   `json:"last_seen_lon"`
}

// NewTiger creates a new Tiger instance.
func NewTiger(name string, dateOfBirth, lastSeenTimestamp time.Time, lastSeenLat, lastSeenLon float64) *Tiger {
	return &Tiger{
		Name:              name,
		DateOfBirth:       dateOfBirth,
		LastSeenTimestamp: lastSeenTimestamp,
		LastSeenLat:       lastSeenLat,
		LastSeenLon:       lastSeenLon,
	}
}

// Save inserts the Tiger into the database.
func (t *Tiger) Save(db *sql.DB) error {
	query := `INSERT INTO tigers (name, date_of_birth, last_seen_timestamp, last_seen_lat, last_seen_lon) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return db.QueryRow(query, t.Name, t.DateOfBirth, t.LastSeenTimestamp, t.LastSeenLat, t.LastSeenLon).Scan(&t.ID)
}

// UpdateLastSeen updates the last seen details of the tiger in the database.
func (t *Tiger) UpdateLastSeen(db *sql.DB, timestamp time.Time, lat, lon float64) error {
	t.LastSeenTimestamp = timestamp
	t.LastSeenLat = lat
	t.LastSeenLon = lon
	query := `UPDATE tigers SET last_seen_timestamp = $2, last_seen_lat = $3, last_seen_lon = $4 WHERE id = $1`
	_, err := db.Exec(query, t.ID, t.LastSeenTimestamp, t.LastSeenLat, t.LastSeenLon)
	return err
}

// GetAllTigers retrieves all tigers from the database with pagination.
func GetAllTigers(db *sql.DB, limit, offset int) ([]Tiger, error) {
	query := `SELECT id, name, date_of_birth, last_seen_timestamp, last_seen_lat, last_seen_lon FROM tigers ORDER BY last_seen_timestamp DESC LIMIT $1 OFFSET $2`
	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tigers []Tiger
	for rows.Next() {
		var t Tiger
		if err := rows.Scan(&t.ID, &t.Name, &t.DateOfBirth, &t.LastSeenTimestamp, &t.LastSeenLat, &t.LastSeenLon); err != nil {
			return nil, err
		}
		tigers = append(tigers, t)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return tigers, nil
}

// GetTigerByID retrieves a single tiger record by its ID from the database.
func GetTigerByID(db *sql.DB, id int) (*Tiger, error) {
	query := `SELECT id, name, date_of_birth, last_seen_timestamp, last_seen_lat, last_seen_lon FROM tigers WHERE id = $1`
	t := Tiger{}
	err := db.QueryRow(query, id).Scan(&t.ID, &t.Name, &t.DateOfBirth, &t.LastSeenTimestamp, &t.LastSeenLat, &t.LastSeenLon)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
