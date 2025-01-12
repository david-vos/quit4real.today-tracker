package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"project/models"
	"project/repository"
)

// AddUserHandler handles adding a new user
func AddUserHandler(db repository.DBExecutor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		var user models.User
		// Parse JSON request body
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Add the user to the database
		err = repository.CreateUser(db, user)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			log.Printf("Error creating user: %v", err)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully"))
	}
}

func AddTrackerHandler(db repository.DBExecutor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, ok := vars["userID"]
		gameId, ok := vars["gameID"]
		if !ok {
			http.Error(w, "userID and GameId are required", http.StatusBadRequest)
			return
		}

		err := repository.CreateTracker(db, userID, gameId)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully"))

	}
}

// TrackUserHandler handles requests for tracking a user by ID
func TrackUserHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, ok := vars["userID"]
		gameId, ok := vars["gameID"]
		if !ok {
			http.Error(w, "userID and GameId are required", http.StatusBadRequest)
			return
		}

		// Retrieve the user from the database
		tracker, err := repository.GetTracker(db, userID, gameId)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find user with ID %d: %v", tracker, err), http.StatusNotFound)
			return
		}

		// Return user details as JSON
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(tracker)
		err = json.NewEncoder(w).Encode(tracker)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		}
	}
}
