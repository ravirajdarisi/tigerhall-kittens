package db

import (
    "database/sql"
    "fmt"
    "log"
    _ "github.com/lib/pq" // PostgreSQL driver
)

// DBConfig holds the database configuration parameters
type DBConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    DBName   string
    SSLMode  string
}

// Connect establishes a connection to the database using the provided configuration.
func Connect(cfg DBConfig) (*sql.DB, error) {
    // Construct the connection string
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

    // Open the connection
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        log.Fatalf("Could not connect to the database: %v", err)
        return nil, err
    }

    // Verify the connection
    err = db.Ping()
    if err != nil {
        log.Fatalf("Failed to ping the database: %v", err)
        return nil, err
    }

    fmt.Println("Successfully connected to the database")
    return db, nil
}
