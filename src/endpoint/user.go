package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/model"
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
		logger.Info("Got request to add a new user")
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
		logger.Info("Got request to add a new tracker")
		vars := mux.Vars(r)
		userID, ok := vars["userID"]
		gameId, ok := vars["gameID"]
		if !ok {
			http.Error(w, "userID and GameId are required", http.StatusBadRequest)
			return
		}

		err := endpoint.TrackerCommandHandler.Add(userID, gameId)
		if err != nil {
			logger.Debug("Failed to add Tracker: " + err.Error())
			http.Error(w, "Failed to add tracker", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		_, err = w.Write([]byte("Tracker created successfully"))
		if err != nil {
			return
		}
	}
}
