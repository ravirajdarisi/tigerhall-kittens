package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/ravirajdarisi/tigerhall-kittens/models"
)


type UserRegistrationRequest struct {
    Username string `json:"username"`
    Email    string `json:"email"`
    Password string `json:"password"` // This is the plaintext password sent by the client
}



// CreateUserHandler handles the creation of a new user.
func CreateUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Decode the request body into a User struct.
		
		var req UserRegistrationRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, "Invalid user data", http.StatusBadRequest)
			return
		}

		// Validate the input for correctness.
		if req.Username == "" || req.Email == "" || req.Password == "" {
			http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
			return
		}

		// Check if the username or email already exists.
		existingUser, _ := models.GetUserByUsername(db, req.Username)
		if existingUser != nil {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		// Create a new User instance and hash the password.
		newUser, err := models.NewUser(req.Username, req.Password, req.Email)
		if err != nil {
			http.Error(w, "Error creating user", http.StatusInternalServerError)
			return
		}

		// Insert the new user record into the database.
		err = newUser.Save(db)
		if err != nil {
			http.Error(w, "Error saving user to the database", http.StatusInternalServerError)
			return
		}

		// Respond to the request indicating the user was created.
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newUser)
	}
}


// LoginHandler handles the user login.
func LoginHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Decode the request body to get the credentials.
		var creds struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Fetch the user from the database.
		user, err := models.GetUserByUsername(db, creds.Username)
		if err != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		// Authenticate the user.
		if !user.Authenticate(db, creds.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// Respond to the request indicating the user was authenticated.
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(struct{ Status string }{"Logged in successfully"})
	}
}





