package models

import (
	"database/sql"
	"time"
	"golang.org/x/crypto/bcrypt"
)

// User represents the user structure.
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"created_at"`
}

// NewUser creates a new User instance and hashes the password.
func NewUser(username, password, email string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		CreatedAt:    time.Now(),
	}, nil
}

// Save inserts the User into the database.
func (u *User) Save(db *sql.DB) error {
	query := `INSERT INTO users (username, password_hash, email, created_at) VALUES ($1, $2, $3, $4)`
	_, err := db.Exec(query, u.Username, u.PasswordHash, u.Email, u.CreatedAt)
	return err
}

// Authenticate checks if the provided password is correct.
func (u *User) Authenticate(db *sql.DB, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

// GetUserByUsername fetches the user with the given username from the database.
func GetUserByUsername(db *sql.DB, username string) (*User, error) {
	query := `SELECT id, username, password_hash, email, created_at FROM users WHERE username = $1`
	user := User{}
	err := db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.PasswordHash, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
