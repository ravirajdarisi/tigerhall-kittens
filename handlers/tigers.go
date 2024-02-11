package handlers

import (

	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ravirajdarisi/tigerhall-kittens/models"
)



// CreateTigerHandler handles the creation of a new tiger.
func CreateTigerHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		

		// Decode the request body into a Tiger struct.
		var newTiger models.Tiger
		err := json.NewDecoder(r.Body).Decode(&newTiger)
		if err != nil {
			http.Error(w, "Invalid tiger data", http.StatusBadRequest)
			return
		}

		// Validate the input data.
		if newTiger.Name == "" || newTiger.DateOfBirth.IsZero() ||
		   newTiger.LastSeenTimestamp.IsZero() ||
		   newTiger.LastSeenLat < -90 || newTiger.LastSeenLat > 90 || newTiger.LastSeenLat == 0 ||
		   newTiger.LastSeenLon < -180 || newTiger.LastSeenLon > 180 || newTiger.LastSeenLon == 0 {
			http.Error(w, "Missing or invalid tiger details: Ensure all fields are present and latitude/longitude are within valid ranges", http.StatusBadRequest)
			return
		}

		// Insert the new tiger record into the database.
		err = newTiger.Save(db)
		if err != nil {
			http.Error(w, "Error saving tiger to the database", http.StatusInternalServerError)
			return
		}

		// Respond to the request indicating the tiger was created.
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newTiger)
	}
}


func ListAllTigersHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse pagination parameters
		page, err := strconv.Atoi(r.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1 // Default to the first page if not specified or invalid
		}

		pageSize, err := strconv.Atoi(r.URL.Query().Get("pageSize"))
		if err != nil || pageSize <= 0 {
			pageSize = 10 // Default page size or enforce a maximum to prevent overly large requests
		}

		// Calculate offset
		offset := (page - 1) * pageSize

		// Retrieve paginated tigers from the database.
		tigers, err := models.GetAllTigers(db, pageSize, offset)
		if err != nil {
			http.Error(w, "Error retrieving tigers from the database", http.StatusInternalServerError)
			return
		}

		// Respond to the request with the list of tigers.
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tigers)
	}
}
