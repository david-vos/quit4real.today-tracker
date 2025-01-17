package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"project/logger"
	"project/main/handler/command"
	"project/main/model"
)

type UserEndpoint struct {
	Router                *mux.Router
	UserCommandHandler    *command.UserCommandHandler
	TrackerCommandHandler *command.TrackerCommandHandler
}

func (endpoint *UserEndpoint) User() {
	logger.Info("Trying to start the user endpoints")
	endpoint.Router.HandleFunc("/users", endpoint.AddUser()).Methods("POST")
	endpoint.Router.HandleFunc("/user/{userID}/track/{gameID}", endpoint.AddTracker()).Methods("POST")
	logger.Info("User endpoints started")
}

func (endpoint *UserEndpoint) AddUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			logger.Debug("Not a POST request")
			http.Error(w, "Only POST method is allowed",
				http.StatusMethodNotAllowed)
			return
		}

		// Parse JSON request body
		var user model.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			logger.Debug("Error decoding user: " + err.Error())
			http.Error(w, "Invalid request body",
				http.StatusBadRequest)
			return
		}

		// Add the user to the database
		var errAddUser = endpoint.UserCommandHandler.Add(user)
		if errAddUser != nil {
			logger.Debug("Error adding user: " + err.Error())
		}

		// Respond with success
		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("User created successfully"))
		if err != nil {
			return
		}
	}
}

func (endpoint *UserEndpoint) AddTracker() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, ok := vars["userID"]
		gameId, ok := vars["gameID"]
		if !ok {
			http.Error(w, "userID and GameId are required", http.StatusBadRequest)
			return
		}

		// this should query the steam API to set the tracker to  a value instead of 0
		err := endpoint.TrackerCommandHandler.Add(userID, gameId)
		if err != nil {
			logger.Debug("Failed to create user" + err.Error())
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("User created successfully"))
		if err != nil {
			return
		}
	}
}
