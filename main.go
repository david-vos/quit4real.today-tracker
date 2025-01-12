package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"project/config"
	"project/db"
	"project/handlers"
)

func main() {
	config.Init()

	dbc := db.Setup()
	defer func() {
		err := dbc.Close()
		if err != nil {
			log.Fatalf("Failed to close the database: %v", err)
		}
	}()

	handlers.SetupCronJobs(dbc)

	router := mux.NewRouter()

	// Define routes
	router.HandleFunc("/users", handlers.AddUserHandler(dbc)).Methods("POST")
	router.HandleFunc("/user/{userID}/track/{gameID}", handlers.TrackUserHandler(dbc)).Methods("GET")
	router.HandleFunc("/user/{userID}/track/{gameID}", handlers.AddTrackerHandler(dbc)).Methods("POST")

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
