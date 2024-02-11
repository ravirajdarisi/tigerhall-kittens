package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/ravirajdarisi/tigerhall-kittens/models"
	"github.com/ravirajdarisi/tigerhall-kittens/utils"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type NotificationMessage struct {
	TigerID int
	UserIDs []int // Users to be notified
}

type SightingRepository interface {
	GetLastSightingByTigerID(tigerID int) (*models.Sighting, error)
	UpdateTigerLastSeen(tigerID int, timestamp time.Time, lat, lon float64) error
	SaveSighting(sighting models.Sighting) error
	GetUsersByTigerID(tigerID int) ([]int, error)
}

type DBSightingRepository struct {
	db *sql.DB
}

func NewDBSightingRepository(db *sql.DB) *DBSightingRepository {
	return &DBSightingRepository{db: db}
}

// GetLastSightingByTigerID retrieves the most recent sighting of a given tiger.
func (repo *DBSightingRepository) GetLastSightingByTigerID(tigerID int) (*models.Sighting, error) {
	sighting := &models.Sighting{}
	query := `SELECT id, tiger_id, lat, lon, timestamp, image_path FROM sightings WHERE tiger_id = $1 ORDER BY timestamp DESC LIMIT 1`
	err := repo.db.QueryRow(query, tigerID).Scan(&sighting.ID, &sighting.TigerID, &sighting.Lat, &sighting.Lon, &sighting.Timestamp, &sighting.ImagePath)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // No sighting found is not an error
		}
		return nil, err
	}
	return sighting, nil
}

func (repo *DBSightingRepository) UpdateTigerLastSeen(tigerID int, timestamp time.Time, lat, lon float64) error {
	query := `UPDATE tigers SET last_seen_timestamp = $2, last_seen_lat = $3, last_seen_lon = $4 WHERE id = $1`
	_, err := repo.db.Exec(query, tigerID, timestamp, lat, lon)
	return err
}

// SaveSighting saves a new sighting to the database.
func (repo *DBSightingRepository) SaveSighting(sighting models.Sighting) error {
	query := `INSERT INTO sightings (user_id, tiger_id, lat, lon, timestamp, image_path) VALUES ($1, $2, $3, $4, $5,$6)`
	_, err := repo.db.Exec(query, sighting.UserID, sighting.TigerID, sighting.Lat, sighting.Lon, sighting.Timestamp, sighting.ImagePath)
	return err
}

// GetUsersByTigerID retrieves a list of unique user IDs who have reported a sighting of a specific tiger.
func (repo *DBSightingRepository) GetUsersByTigerID(tigerID int) ([]int, error) {
	var userIDs []int
	query := `SELECT DISTINCT user_id FROM sightings WHERE tiger_id = $1`
	rows, err := repo.db.Query(query, tigerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var userID int
		if err := rows.Scan(&userID); err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return userIDs, nil
}


func CreateSightingHandler(repo SightingRepository, notificationQueue chan NotificationMessage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			http.Error(w, "Failed to parse multipart form: "+err.Error(), http.StatusBadRequest)
			return
		}

		sightingInfo := r.FormValue("sightingInfo")
		var newSighting models.Sighting
		if err := json.Unmarshal([]byte(sightingInfo), &newSighting); err != nil {
			http.Error(w, "Invalid sighting data", http.StatusBadRequest)
			return
		}

		// Perform validations
		if validationErr := validateSighting(newSighting); validationErr != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationErr)
			return
		}

		// Process the image upload synchronously
		file, header, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "Could not get uploaded file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		imgData, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "Failed to read uploaded file", http.StatusInternalServerError)
			return
		}

		ext := filepath.Ext(header.Filename)
		imagePath, err := processImageUpload(imgData, ext)
		if err != nil {
			log.Printf("Failed to process image upload: %v", err)
			http.Error(w, "Failed to process image upload", http.StatusInternalServerError)
			return
		}
		newSighting.ImagePath = imagePath

		// Database operations
		lastSighting, err := repo.GetLastSightingByTigerID(newSighting.TigerID)

		log.Println(newSighting.UserID)
		if err != nil {
			http.Error(w, "Error retrieving last sighting", http.StatusInternalServerError)
			return
		}

		var filteredUserIDs []int

		if lastSighting != nil {
			// Log the latitudes and longitudes being compared
			log.Printf("Last sighting lat: %v, lon: %v", lastSighting.Lat, lastSighting.Lon)
			log.Printf("New sighting lat: %v, lon: %v", newSighting.Lat, newSighting.Lon)

			// Check distance from the last sighting
			distance := utils.CalculateDistance(lastSighting.Lat, lastSighting.Lon, newSighting.Lat, newSighting.Lon)

			// Log the calculated distance
			log.Printf("Calculated distance: %v kilometers", distance)

			if distance < 5 {
				errMsg := ErrorResponse{
					Code:    "TOO_CLOSE_TO_PREVIOUS_SIGHTING",
					Message: "New sighting is too close to the last sighting. Sightings must be at least 5 kilometers apart.",
				}
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(errMsg)
				return
			}


			// Fetch user IDs of those who have previously sighted the same tiger
			userIDs, err := repo.GetUsersByTigerID(newSighting.TigerID)
			if err != nil {
				http.Error(w, "Failed to fetch users who have sighted the tiger", http.StatusInternalServerError)
				return
			}

			// Filter out the current user from the notifcation list
			for _, userID := range userIDs {
				if userID != newSighting.UserID {
					filteredUserIDs = append(filteredUserIDs, userID)
				}
			}
		}

		
		if err := repo.UpdateTigerLastSeen(newSighting.TigerID, newSighting.Timestamp, newSighting.Lat, newSighting.Lon); err != nil {
			http.Error(w, "Failed to update tiger's last seen information", http.StatusInternalServerError)
			return
		}

		if err := repo.SaveSighting(newSighting); err != nil {
			http.Error(w, "Failed to save sighting", http.StatusInternalServerError)
			return
		}

		// Send notifications if there are users to notify
		if len(filteredUserIDs) > 0 {
			go func(userIDs []int, tigerID int) {
				log.Print("Inside notification block")
				log.Printf("Queuing notifications for User IDs: %v, for Tiger ID: %d", userIDs, tigerID)
				notificationMessage := NotificationMessage{UserIDs: userIDs, TigerID: tigerID}
				notificationQueue <- notificationMessage
			}(filteredUserIDs, newSighting.TigerID)
		}
		

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newSighting)
	}
}

func validateSighting(newSighting models.Sighting) *ErrorResponse {
	// Validate Timestamp
	if newSighting.Timestamp.IsZero() {
		return &ErrorResponse{
			Code:    "INVALID_TIMESTAMP",
			Message: "Timestamp is required and must be a valid date.",
		}
	}

	// Validate Latitude
	if newSighting.Lat < -90 || newSighting.Lat > 90 || newSighting.Lat == 0 {
		return &ErrorResponse{
			Code:    "INVALID_LATITUDE",
			Message: "Latitude must be between -90 and 90 and not zero.",
		}
	}

	// Validate Longitude
	if newSighting.Lon < -180 || newSighting.Lon > 180 || newSighting.Lon == 0 {
		return &ErrorResponse{
			Code:    "INVALID_LONGITUDE",
			Message: "Longitude must be between -180 and 180 and not zero.",
		}
	}

	// Validate UserID
	if newSighting.UserID <= 0 {
		return &ErrorResponse{
			Code:    "INVALID_USER_ID",
			Message: "UserID is invalid.",
		}
	}

	// Validate TigerID
	if newSighting.TigerID <= 0 {
		return &ErrorResponse{
			Code:    "INVALID_TIGER_ID",
			Message: "TigerID is invalid.",
		}
	}

	// All validations passed
	return nil
}

func processImageUpload(imgData []byte, ext string) (string, error) {
	// Determine the storage path.
	storagePath := os.Getenv("IMAGE_STORAGE_PATH")
	if storagePath == "" {
		storagePath = "./default_storage" // Default path relative to the project
	}

	// Ensure the storage directory exists.
	if _, err := os.Stat(storagePath); os.IsNotExist(err) {
		err := os.MkdirAll(storagePath, 0755) // Create the directory with read & execute permissions
		if err != nil {
			return "", fmt.Errorf("failed to create storage directory: %v", err)
		}
	}

	// Decode the image from the buffered data.
	var img image.Image
	var err error
	switch ext {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(bytes.NewReader(imgData))
	case ".png":
		img, err = png.Decode(bytes.NewReader(imgData))
	default:
		return "", fmt.Errorf("unsupported file type")
	}
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	resizedImg := utils.ResizeImage(img, 250, 250)
	fileName := "resized_image" + ext
	filePath := filepath.Join(storagePath, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %v", err)
	}
	defer dst.Close()
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(dst, resizedImg, nil)
	case ".png":
		err = png.Encode(dst, resizedImg)
	}
	if err != nil {
		return "", fmt.Errorf("failed to write image to file: %v", err)
	}

	return filePath, nil
}

// ListSightingsHandler creates an HTTP handler function for listing sightings with pagination.
func ListSightingsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
	
		tigerIDStr := r.URL.Query().Get("tigerID")
		if tigerIDStr == "" {
			http.Error(w, "Tiger ID is required", http.StatusBadRequest)
			return
		}
		tigerID, err := strconv.Atoi(tigerIDStr)
		if err != nil {
			http.Error(w, "Invalid Tiger ID", http.StatusBadRequest)
			return
		}

		// Parse pagination parameters
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1 
		}

		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil || pageSize <= 0 {
			pageSize = 10
		}

		// Calculate offset
		offset := (page - 1) * pageSize

		// Fetch sightings with pagination
		sightings, err := models.GetAllSightingsByTigerID(db, tigerID, pageSize, offset)
		if err != nil {
			http.Error(w, "Failed to fetch sightings", http.StatusInternalServerError)
			return
		}

		// Respond with the list of sightings in JSON format
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(sightings)
	}
}
