package main

import (
	"log"
	"net/http"

	"event-registration-system/internal/handlers"
	"event-registration-system/internal/repository"
	"event-registration-system/pkg/database"
)

func main() {
	// Initialize database
	if err := database.InitDB("event_registration.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Initialize repos and handlers
	db := database.DB

	eventRepo := repository.NewEventRepository(db)
	regRepo := repository.NewRegistrationRepository(db)
	userRepo := repository.NewUserRepository(db)

	eventHandler := handlers.NewEventHandler(eventRepo)
	regHandler := handlers.NewRegistrationHandler(regRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	// Set up router using standard library
	mux := http.NewServeMux()

	// Organizer routes
	mux.HandleFunc("/api/events/create", eventHandler.CreateEvent)

	// User routes
	mux.HandleFunc("/api/events", eventHandler.GetEvents)
	mux.HandleFunc("/api/register", regHandler.RegisterUser)
	mux.HandleFunc("/api/users/create", userHandler.CreateUser)

	log.Println("Server executing on :8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
