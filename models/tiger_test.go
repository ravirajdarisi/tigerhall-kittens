package models

import (
	"testing"
	"time"
    _"database/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/DATA-DOG/go-sqlmock"
)

func TestNewTiger(t *testing.T) {
	name := "TigerName"
	dateOfBirth := time.Now()
	lastSeenTimestamp := time.Now()
	lastSeenLat := 10.123
	lastSeenLon := 20.456

	tiger := NewTiger(name, dateOfBirth, lastSeenTimestamp, lastSeenLat, lastSeenLon)

	assert.NotNil(t, tiger, "NewTiger should return a non-nil instance of Tiger")
	assert.Equal(t, name, tiger.Name, "Tiger name should be the same as what was passed into NewTiger")
	assert.Equal(t, dateOfBirth, tiger.DateOfBirth, "Tiger date of birth should be the same as what was passed into NewTiger")
	assert.Equal(t, lastSeenTimestamp, tiger.LastSeenTimestamp, "Tiger last seen timestamp should be the same as what was passed into NewTiger")
	assert.Equal(t, lastSeenLat, tiger.LastSeenLat, "Tiger last seen latitude should be the same as what was passed into NewTiger")
	assert.Equal(t, lastSeenLon, tiger.LastSeenLon, "Tiger last seen longitude should be the same as what was passed into NewTiger")
}



func TestTiger_Save(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	tiger := &Tiger{
		Name:              "TigerTest",
		DateOfBirth:       time.Now(),
		LastSeenTimestamp: time.Now(),
		LastSeenLat:       12.34,
		LastSeenLon:       56.78,
	}

	mock.ExpectQuery("INSERT INTO tigers").
		WithArgs(tiger.Name, tiger.DateOfBirth, tiger.LastSeenTimestamp, tiger.LastSeenLat, tiger.LastSeenLon).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err = tiger.Save(db)
	assert.NoError(t, err)
	assert.Equal(t, 1, tiger.ID, "After saving, tiger ID should be set to 1")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}



func TestTiger_UpdateLastSeen(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Setting up the tiger instance
	tiger := &Tiger{
		ID: 1, // Assuming the tiger with ID 1 exists in the database
	}

	// The new last seen details to be updated
	newTimestamp := time.Now()
	newLat := 12.345
	newLon := 67.890

	// Setting up the expected database operation
	mock.ExpectExec("UPDATE tigers SET").
		WithArgs(tiger.ID, newTimestamp, newLat, newLon).
		WillReturnResult(sqlmock.NewResult(0, 1)) // Assuming one row is affected

	// Calling the method under test
	err = tiger.UpdateLastSeen(db, newTimestamp, newLat, newLon)
	assert.NoError(t, err)

	// Verifying the tiger struct is updated
	assert.Equal(t, newTimestamp, tiger.LastSeenTimestamp, "Tiger's last seen timestamp should be updated")
	assert.Equal(t, newLat, tiger.LastSeenLat, "Tiger's last seen latitude should be updated")
	assert.Equal(t, newLon, tiger.LastSeenLon, "Tiger's last seen longitude should be updated")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}



func TestGetTigerByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Example tiger data to be returned
	tigerID := 1
	tigerRow := sqlmock.NewRows([]string{"id", "name", "date_of_birth", "last_seen_timestamp", "last_seen_lat", "last_seen_lon"}).
		AddRow(tigerID, "TigerOne", time.Now(), time.Now(), 10.123, 20.123)

	// Setting up the expected query for a specific tiger ID
	mock.ExpectQuery("SELECT id, name, date_of_birth, last_seen_timestamp, last_seen_lat, last_seen_lon FROM tigers WHERE id =").
		WithArgs(tigerID).
		WillReturnRows(tigerRow)

	// Calling the method under test
	tiger, err := GetTigerByID(db, tigerID)
	require.NoError(t, err)

	// Asserting the expected outcomes
	assert.NotNil(t, tiger, "Expected a non-nil tiger to be returned")
	assert.Equal(t, tigerID, tiger.ID, "The tiger's ID should match the queried ID")
	assert.Equal(t, "TigerOne", tiger.Name, "The tiger's name should match")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}



func TestGetAllTigers(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Example data to be returned
	tigerRows := sqlmock.NewRows([]string{"id", "name", "date_of_birth", "last_seen_timestamp", "last_seen_lat", "last_seen_lon"}).
		AddRow(1, "TigerOne", time.Now(), time.Now(), 10.123, 20.123).
		AddRow(2, "TigerTwo", time.Now(), time.Now(), 15.123, 25.123)

	// Setting up the expected query with pagination parameters
	limit := 2
	offset := 0
	mock.ExpectQuery("SELECT id, name, date_of_birth, last_seen_timestamp, last_seen_lat, last_seen_lon FROM tigers ORDER BY last_seen_timestamp DESC").
		WithArgs(limit, offset).
		WillReturnRows(tigerRows)

	// Calling the method under test
	tigers, err := GetAllTigers(db, limit, offset)
	require.NoError(t, err)

	// Asserting the expected outcomes
	require.Len(t, tigers, 2, "Expected two tigers to be returned")
	assert.Equal(t, 1, tigers[0].ID, "The first tiger's ID should match")
	assert.Equal(t, "TigerOne", tigers[0].Name, "The first tiger's name should match")
	assert.Equal(t, 2, tigers[1].ID, "The second tiger's ID should match")
	assert.Equal(t, "TigerTwo", tigers[1].Name, "The second tiger's name should match")

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}




