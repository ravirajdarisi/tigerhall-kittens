package models

import (
	_ "database/sql"
	"reflect"
	"testing"
	"time"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestSighting_Save(t *testing.T) {
    // Create a new mock SQL database connection.
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    // Prepare a mock expectation for the SQL INSERT statement.
    mock.ExpectExec("INSERT INTO sightings").
        WithArgs(1, 50.0, 30.0, sqlmock.AnyArg(), "path/to/image").
        WillReturnResult(sqlmock.NewResult(1, 1)) // Mocking that 1 row was affected.

	
    // Create a new Sighting instance with test data.
    sighting := Sighting{
        UserID:    1, 
        TigerID:   1,
        Lat:       50.0,
        Lon:       30.0,
        Timestamp: time.Now(),
        ImagePath: "path/to/image",
    }

    // Attempt to save the sighting using the mock database connection.
    err = sighting.Save(db)
    require.NoError(t, err)

    // Ensure all expectations set on the mock database were met.
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("there were unfulfilled expectations: %s", err)
    }
}


func TestGetLastSightingByTigerID(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    // Mock data
    rows := sqlmock.NewRows([]string{"id", "tiger_id", "lat", "lon", "timestamp", "image_path"}).
        AddRow(1, 1, 10.1234, 20.5678, time.Now(), "/path/to/image.jpg")

    // Setting up the expectation
    mock.ExpectQuery("SELECT id, tiger_id, lat, lon, timestamp, image_path FROM sightings WHERE tiger_id = \\$1 ORDER BY timestamp DESC LIMIT 1").
        WithArgs(1).
        WillReturnRows(rows)

    // Calling the method under test
    sighting, err := GetLastSightingByTigerID(db, 1)
    require.NoError(t, err)
    require.NotNil(t, sighting)
    require.Equal(t, 1, sighting.TigerID)
    require.Equal(t, "/path/to/image.jpg", sighting.ImagePath)

    // Ensure all expectations were met
    err = mock.ExpectationsWereMet()
    require.NoError(t, err)
}


func TestGetAllSightingsByTigerID(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    // Mock data
    rows := sqlmock.NewRows([]string{"id", "tiger_id", "lat", "lon", "timestamp", "image_path"}).
        AddRow(1, 1, 10.1234, 20.5678, time.Now(), "/images/image1.jpg").
        AddRow(2, 1, 11.1234, 21.5678, time.Now(), "/images/image2.jpg")

    // Setting up the expectation
    tigerID, limit, offset := 1, 2, 0
    mock.ExpectQuery("SELECT id, tiger_id, lat, lon, timestamp, image_path FROM sightings WHERE tiger_id = \\$1 ORDER BY timestamp DESC LIMIT \\$2 OFFSET \\$3").
        WithArgs(tigerID, limit, offset).
        WillReturnRows(rows)

    // Calling the method under test
    sightings, err := GetAllSightingsByTigerID(db, tigerID, limit, offset)
    require.NoError(t, err)
    require.Len(t, sightings, 2)
    require.True(t, reflect.DeepEqual(sightings[0].TigerID, tigerID))
    require.Equal(t, "/images/image1.jpg", sightings[0].ImagePath)
    require.Equal(t, "/images/image2.jpg", sightings[1].ImagePath)

    // Ensure all expectations were met
    err = mock.ExpectationsWereMet()
    require.NoError(t, err)
}


func TestGetUsersByTigerID(t *testing.T) {
    db, mock, err := sqlmock.New()
    require.NoError(t, err)
    defer db.Close()

    // Prepare mock data
    tigerID := 1
    rows := sqlmock.NewRows([]string{"user_id"}).
        AddRow(101).
        AddRow(102).
        AddRow(103)

    // Setting up the expectation
    mock.ExpectQuery("SELECT DISTINCT user_id FROM sightings WHERE tiger_id = \\$1").
        WithArgs(tigerID).
        WillReturnRows(rows)

    // Calling the function under test
    userIDs, err := GetUsersByTigerID(db, tigerID)
    require.NoError(t, err)
    require.NotNil(t, userIDs)
    require.Len(t, userIDs, 3)
    require.Equal(t, 101, userIDs[0])
    require.Equal(t, 102, userIDs[1])
    require.Equal(t, 103, userIDs[2])

    // Ensure all expectations were met
    err = mock.ExpectationsWereMet()
    require.NoError(t, err)
}