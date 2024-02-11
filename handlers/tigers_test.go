package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ravirajdarisi/tigerhall-kittens/models"
	"github.com/stretchr/testify/assert"
)

// TestCreateTigerHandler tests various scenarios for the CreateTigerHandler function.
func TestCreateTigerHandler(t *testing.T) {
	
	db, mock := setupMockDB(t) 

	mock.ExpectQuery("INSERT INTO tigers").
    WithArgs(
        sqlmock.AnyArg(), // Name
        sqlmock.AnyArg(), // DateOfBirth
        sqlmock.AnyArg(), // LastSeenTimestamp
        sqlmock.AnyArg(), // LastSeenLat
        sqlmock.AnyArg() , // LastSeenLon
    ).WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1)) 
	handler := CreateTigerHandler(db)

	tests := []struct {
		name           string
		tiger          models.Tiger
		expectedStatus int
	}{
		{
			name: "Valid request",
			tiger: models.Tiger{
				Name:              "ValidTiger",
				DateOfBirth:       time.Now(),
				LastSeenTimestamp: time.Now(),
				LastSeenLat:       45.0,
				LastSeenLon:       90.0,
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "Invalid latitude",
			tiger: models.Tiger{
				Name:              "TigerWithBadLat",
				DateOfBirth:       time.Now(),
				LastSeenTimestamp: time.Now(),
				LastSeenLat:       -91.0, 
				LastSeenLon:       90.0,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid longitude",
			tiger: models.Tiger{
				Name:              "TigerWithBadLon",
				DateOfBirth:       time.Now(),
				LastSeenTimestamp: time.Now(),
				LastSeenLat:       45.0,
				LastSeenLon:       181.0, 
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal the tiger struct to JSON for the request body
			body, err := json.Marshal(tt.tiger)
			assert.NoError(t, err)
			req, err := http.NewRequest("POST", "/tiger", bytes.NewBuffer(body))
			assert.NoError(t, err)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)
			assert.Equal(t, tt.expectedStatus, rr.Code, "handler returned wrong status code")
		})
	}
}




func TestListAllTigersHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	layout := "2006-01-02"
	dateOfBirthOne, _ := time.Parse(layout, "2020-01-01")
	lastSeenTimestampOne, _ := time.Parse(layout, "2021-01-01")
	dateOfBirthTwo, _ := time.Parse(layout, "2020-02-02")
	lastSeenTimestampTwo, _ := time.Parse(layout, "2021-02-02")

	// Mock the GetAllTigers query
	rows := sqlmock.NewRows([]string{"id", "name", "date_of_birth", "last_seen_timestamp", "last_seen_lat", "last_seen_lon"}).
		AddRow(1, "TigerOne", dateOfBirthOne, lastSeenTimestampOne, 10.123, 20.123).
		AddRow(2, "TigerTwo", dateOfBirthTwo, lastSeenTimestampTwo, 30.123, 40.123)
	mock.ExpectQuery("^SELECT id, name, date_of_birth, last_seen_timestamp, last_seen_lat, last_seen_lon FROM tigers").
		WillReturnRows(rows)

	handler := ListAllTigersHandler(db)

	// Create a request to pass to our handler.
	req, err := http.NewRequest("GET", "/tigers?page=1&pageSize=2", nil)
	assert.NoError(t, err)

	// Create a ResponseRecorder to record the response.
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	assert.Equal(t, http.StatusOK, rr.Code)
}

func setupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}

	return db, mock
}
