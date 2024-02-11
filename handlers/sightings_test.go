package handlers

import (
	"bytes"
	_"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"github.com/ravirajdarisi/tigerhall-kittens/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSightingRepository is a mock type for the SightingRepository interface
type MockSightingRepository struct {
	mock.Mock
}

// Define methods that match the interface you are mocking
func (m *MockSightingRepository) GetLastSightingByTigerID(tigerID int) (*models.Sighting, error) {
	args := m.Called(tigerID)
	return args.Get(0).(*models.Sighting), args.Error(1)
}

func (m *MockSightingRepository) UpdateTigerLastSeen(tigerID int, timestamp time.Time, lat, lon float64) error {
	args := m.Called(tigerID, timestamp, lat, lon)
	return args.Error(0)
}

func (m *MockSightingRepository) SaveSighting(sighting models.Sighting) error {
	args := m.Called(sighting)
	return args.Error(0)
}

func (m *MockSightingRepository) GetUsersByTigerID(tigerID int) ([]int, error) {
	args := m.Called(tigerID)
	return args.Get(0).([]int), args.Error(1)
}

func TestCreateSightingHandler(t *testing.T) {
	mockRepo := new(MockSightingRepository)
	dummyNotificationQueue := make(chan NotificationMessage, 1)
	handler := CreateSightingHandler(mockRepo,dummyNotificationQueue)

	// Setup mock behavior
	mockSighting := &models.Sighting{} 
	mockRepo.On("GetLastSightingByTigerID", mock.Anything).Return(mockSighting, nil)
	mockRepo.On("UpdateTigerLastSeen", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	mockRepo.On("SaveSighting", mock.Anything).Return(nil)
	mockRepo.On("GetUsersByTigerID", mock.Anything).Return([]int{}, nil)

	
	// Simulate a valid sighting JSON payload
	sighting := models.Sighting{
		UserID:    1,
		TigerID:   1,
		Lat:       10.0,
		Lon:       20.0,
		Timestamp: time.Now(),
		ImagePath: "path/to/dummy/image.jpg", // Added dummy image path
	}

	payload, _ := json.Marshal(sighting)

	req, _ := http.NewRequest("POST", "/sighting", bytes.NewBuffer(payload))
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	responseBody := rr.Body.String()
	fmt.Println("Response body:", responseBody)
	assert.Equal(t, http.StatusCreated, rr.Code)
	mockRepo.AssertExpectations(t)
}




