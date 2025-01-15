package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"project/config"
	"project/models"
	"project/repository"
)

type UserController struct {
	UserRepoContr    *repository.UserRepoController
	TrackerRepoContr *repository.TrackerRepoController
}

// AddUserHandler handles adding a new user in the DataBase
func (c *UserController) AddUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed",
				http.StatusMethodNotAllowed)
			return
		}

		var user models.User
		// Parse JSON request body
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			config.HandleError("Invalid request body", err)
			http.Error(w, "Invalid request body",
				http.StatusBadRequest)
			return
		}

		// Add the user to the database
		err = c.UserRepoContr.CreateUser(user)
		if err != nil {
			http.Error(w, "Failed to create user",
				http.StatusInternalServerError)
			config.HandleError("Error creating user: %v", err)
			return
		}

		// Respond with success
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully"))
	}
}

func (c *UserController) AddTrackerHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, ok := vars["userID"]
		gameId, ok := vars["gameID"]
		if !ok {
			http.Error(w, "userID and GameId are required", http.StatusBadRequest)
			return
		}

		// this should query the steam API to set the tracker to  a value instead of 0

		err := c.TrackerRepoContr.CreateTracker(userID, gameId)
		if err != nil {
			config.HandleError("Failed to create user", err)
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("User created successfully"))

	}
}

// TrackUserHandler handles requests for tracking a user by ID
func (c *UserController) TrackUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, ok := vars["userID"]
		gameId, ok := vars["gameID"]
		if !ok {
			http.Error(w, "userID and GameId are required", http.StatusBadRequest)
			return
		}

		// Retrieve the user from the database
		tracker, err := c.TrackerRepoContr.GetTracker(userID, gameId)
		if err != nil {
			config.HandleError("Failed to find user", err)
			http.Error(w, fmt.Sprintf("Failed to find user with ID %d: %v", tracker, err), http.StatusNotFound)
			return
		}

		// Return user details as JSON
		w.Header().Set("Content-Type", "application/json")
		fmt.Println(tracker)
		err = json.NewEncoder(w).Encode(tracker)
		if err != nil {
			config.HandleError("Failed to write response", err)
			http.Error(w, fmt.Sprintf("Failed to write response: %v", err), http.StatusInternalServerError)
		}
	}
}
