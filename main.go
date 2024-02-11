package main

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ravirajdarisi/tigerhall-kittens/db"
	"github.com/ravirajdarisi/tigerhall-kittens/handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var notificationQueue chan handlers.NotificationMessage
var wg sync.WaitGroup

func main() {
	// Initialize the notificationQueue
	notificationQueue = make(chan handlers.NotificationMessage, 100)

	// Setup database configuration
	dbConfig := db.DBConfig{
		Host:     "localhost",
		Port:     5432,
		User:     "my-postgres-cluster",
		Password: "mIw1-rc^ZRL2r)V@.c6L9:qW",
		DBName:   "mydb",
		SSLMode:  "disable",
	}

	db, err := db.Connect(dbConfig)
	if err != nil {
		log.Fatalf("Could not connect to the database: %v", err)
	}
	defer db.Close()

	// Context to control cancellation
	// need to improve this logic to pass the context acorss all the handlers
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize HTTP routes
	setupRoutes(db, ctx)

	// Start the notification processor
	wg.Add(1)
	go startNotificationProcessor(ctx)

	// Start HTTP server in a goroutine
	server := &http.Server{Addr: ":8080", Handler: nil}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	fmt.Println("Server started on :8080")

	// Graceful shutdown setup
	gracefulShutdown(ctx, server)
}

func setupRoutes(db *sql.DB, ctx context.Context) {

	sightingRepo := handlers.NewDBSightingRepository(db)
	// list of all handlers
	// use the above ctx to handlers for proper graceful shutdowns
	http.HandleFunc("/users/create", handlers.CreateUserHandler(db))
	http.HandleFunc("/users/login", handlers.LoginHandler(db))
	http.HandleFunc("/tigers/create", handlers.CreateTigerHandler(db))
	http.HandleFunc("/tigers/list", handlers.ListAllTigersHandler(db))
	http.HandleFunc("/sightings/create", handlers.CreateSightingHandler(sightingRepo, notificationQueue))
	http.HandleFunc("/sightings/list", handlers.ListSightingsHandler(db))

}

func gracefulShutdown(ctx context.Context, server *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs // Block until a signal is received
	fmt.Println("Received shutdown signal")

	// Initiate graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		fmt.Printf("HTTP server Shutdown: %v", err)
	}

	close(notificationQueue) // Close the notification channel
	wg.Wait()                // Wait for the notification processor to finish

	fmt.Println("Server shutdown gracefully")
}

func startNotificationProcessor(ctx context.Context) {
	defer wg.Done()

	for {
		select {
		case message, ok := <-notificationQueue:
			if !ok {
				fmt.Println("Notification queue closed, stopping processor")
				return // Exit the loop and goroutine
			}
			log.Printf("Processing notification for User IDs: %v, for Tiger ID: %d", message.UserIDs, message.TigerID)
			for _, userID := range message.UserIDs {
				log.Printf("Sending email to User ID: %d, for Tiger ID: %d", userID, message.TigerID)
				sendEmailToUser(userID, message.TigerID)
			}
		case <-ctx.Done():
			fmt.Println("Shutdown signal received, stopping notification processor")
			return // Exit the loop and goroutine
		}
	}
}


func sendEmailToUser(userID int, tigerID int) {
	// Placeholder for email sending logic
	log.Print("inside main notification in the main block")
	fmt.Printf("Sending email to user %d about tiger %d\n", userID, tigerID)
}
