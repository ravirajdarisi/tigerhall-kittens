package models

import (
	"database/sql"
	_ "database/sql"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"golang.org/x/crypto/bcrypt"
)

func TestNewUser(t *testing.T) {
	username := "testuser"
	password := "password"
	email := "test@example.com"

	user, err := NewUser(username, password, email)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != username {
		t.Errorf("Expected username %v, got %v", username, user.Username)
	}

	if user.Email != email {
		t.Errorf("Expected email %v, got %v", email, user.Email)
	}

	// Check password hash
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		t.Errorf("Password hash does not match the original password")
	}

	if user.CreatedAt.After(time.Now()) {
		t.Errorf("CreatedAt should be before the current time")
	}
}




func TestSave(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	username := "testuser"
	password := "password"
	email := "test@example.com"
	user, _ := NewUser(username, password, email)

	mock.ExpectExec("INSERT INTO users").
		WithArgs(user.Username, user.PasswordHash, user.Email, user.CreatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	if err := user.Save(db); err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}



func TestAuthenticate(t *testing.T) {
	password := "password"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := User{
		Username:     "testuser",
		PasswordHash: string(hashedPassword),
		Email:        "test@example.com",
	}

	// Test with correct password
	if !user.Authenticate(nil, password) {
		t.Errorf("Authenticate failed with the correct password")
	}

	// Test with incorrect password
	if user.Authenticate(nil, "wrongpassword") {
		t.Errorf("Authenticate succeeded with the wrong password")
	}
}



func TestGetUserByUsername(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	username := "testuser"
	rows := sqlmock.NewRows([]string{"id", "username", "password_hash", "email", "created_at"}).
		AddRow(1, username, "hashedpassword", "test@example.com", time.Now())

	mock.ExpectQuery("SELECT id, username, password_hash, email, created_at FROM users WHERE username =").
		WithArgs(username).
		WillReturnRows(rows)

	user, err := GetUserByUsername(db, username)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if user.Username != username {
		t.Errorf("Expected username %v, got %v", username, user.Username)
	}

	// Check if expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}

func TestGetUserByUsername_NotFound(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("SELECT id, username, password_hash, email, created_at FROM users WHERE username =").
		WithArgs("nonexistent").
		WillReturnError(sql.ErrNoRows)

	_, err = GetUserByUsername(db, "nonexistent")
	if err == nil {
		t.Fatalf("Expected an error, got none")
	}

	if !reflect.DeepEqual(err, sql.ErrNoRows) {
		t.Errorf("Expected sql.ErrNoRows, got %v", err)
	}

	// Check if expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("There were unfulfilled expectations: %s", err)
	}
}
