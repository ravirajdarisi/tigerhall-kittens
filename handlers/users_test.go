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
	"golang.org/x/crypto/bcrypt"
)



func TestCreateUserHandler(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, username, password_hash, email, created_at FROM users WHERE username =").
		WithArgs("testuser").
		WillReturnError(sql.ErrNoRows)

	
	mock.ExpectExec("INSERT INTO users").
		WithArgs("testuser", sqlmock.AnyArg(), "test@example.com", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	handler := CreateUserHandler(db)

	
	userData := map[string]string{
		"username": "testuser",
		"password": "password", 
		"email":    "test@example.com",
	}
	userDataJSON, _ := json.Marshal(userData)
	req, _ := http.NewRequest("POST", "/create", bytes.NewBuffer(userDataJSON))
	req.Header.Set("Content-Type", "application/json") // Ensure the request Content-Type is set.
	w := httptest.NewRecorder()

	// Call the handler
	handler.ServeHTTP(w, req)

	// Check the response code
	if status := w.Code; status != http.StatusCreated {
		t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	// Attempt to decode the response body into a User struct.
	var u models.User // Adjust this to the correct import path or definition for User struct.
	err = json.NewDecoder(w.Body).Decode(&u)
	if err != nil {
		t.Errorf("Could not decode response body into user struct: %v", err)
	}

	// Validate the response content. Adjust validations if necessary, based on actual response structure.
	if u.Username != "testuser" || u.Email != "test@example.com" {
		t.Errorf("Handler returned unexpected body: got username %v, email %v", u.Username, u.Email)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}



func TestLoginHandler(t *testing.T) {
    db, mock, err := sqlmock.New()
    if err != nil {
        t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
    }
    defer db.Close()

    // Mock user data
    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
    rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "email", "created_at"}).
        AddRow(1, "testuser", string(hashedPassword), "test@example.com", time.Now())

    // Scenario 1: Successful login
    mock.ExpectQuery("SELECT id, username, password_hash, email, created_at FROM users WHERE username =").
        WithArgs("testuser").
        WillReturnRows(rows)

    // Create the handler
    handler := LoginHandler(db)

    // Create a request to pass to our handler
    credentials := map[string]string{
        "username": "testuser",
        "password": "password",
    }
    credentialsJSON, _ := json.Marshal(credentials)
    req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(credentialsJSON))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    handler.ServeHTTP(w, req)

    // Check the status code is what we expect
    if status := w.Code; status != http.StatusOK {
        t.Errorf("Handler returned wrong status code: got %v want %v", status, http.StatusOK)
    }

    // Scenario 2: Invalid credentials 
    mock.ExpectQuery("SELECT id, username, password_hash, email, created_at FROM users WHERE username =").
        WithArgs("testuser").
        WillReturnRows(rows) // Assuming the password provided doesn't match

    credentialsWrongPassword := map[string]string{
        "username": "testuser",
        "password": "wrongpassword",
    }
    credentialsWrongJSON, _ := json.Marshal(credentialsWrongPassword)
    reqWrong, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(credentialsWrongJSON))
    reqWrong.Header.Set("Content-Type", "application/json")
    wWrong := httptest.NewRecorder()

    handler.ServeHTTP(wWrong, reqWrong)

    // Check the status code is what we expect for wrong credentials
    if status := wWrong.Code; status != http.StatusUnauthorized {
        t.Errorf("Handler returned wrong status code for wrong credentials: got %v want %v", status, http.StatusUnauthorized)
    }

    // Ensure all expectations were met
    if err := mock.ExpectationsWereMet(); err != nil {
        t.Errorf("There were unfulfilled expectations: %s", err)
    }
}



