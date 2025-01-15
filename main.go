package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"project/config"
	"project/db"
)

func main() {

	// The config part
	config.Init()
	dbc := db.Setup()
	defer func() {
		err := dbc.Close()
		if err != nil {
			config.HandleError("Failed to close the database: %v", err)
		}
	}()
	userHandlerContrImp, cronContrImp, failContrImp := DependencyInit(dbc)

	// The cron-job part
	cronContrImp.SetupCronJobs()

	// The web-api part
	router := mux.NewRouter()
	// Define routes
	router.HandleFunc("/users", userHandlerContrImp.AddUserHandler()).Methods("POST")
	router.HandleFunc("/user/{userID}/track/{gameID}", userHandlerContrImp.TrackUserHandler()).Methods("GET")
	router.HandleFunc("/user/{userID}/track/{gameID}", userHandlerContrImp.AddTrackerHandler()).Methods("POST")

	router.HandleFunc("/fail/leaderboard", failContrImp.GetFailsLeaderBoard()).Methods("GET")

	log.Println("Server is running on http://localhost:8080")

	// Start the HTTP server in a separate goroutine
	go func() {
		if err := http.ListenAndServe(":8080", router); err != nil {
			config.HandleError("Failed to start server: %v", err)
		}
	}()

	// Keep the main function running to allow cron jobs to execute
	select {}
}
