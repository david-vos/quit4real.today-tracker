package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"quit4real.today/logger"
	"quit4real.today/src/api"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/model"
)

type UserEndpoint struct {
	Router                     *mux.Router
	SteamApi                   *api.SteamApi
	UserCommandHandler         *command.UserCommandHandler
	SubscriptionCommandHandler *command.SubscriptionCommandHandler
}

func (endpoint *UserEndpoint) User() {
	logger.Info("Trying to start the user endpoints")
	endpoint.Router.HandleFunc("/users", endpoint.AddUser()).Methods("POST")
	endpoint.Router.HandleFunc("/users/${userName}/steamId", endpoint.GetSteamId()).Methods("GET")
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

func (endpoint *UserEndpoint) GetSteamId() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		userID, ok := vars["userName"]
		if !ok {
			http.Error(w, "userName is required", http.StatusBadRequest)
			return
		}

		steamId, err := endpoint.SteamApi.GetSteamIdFromVanityName(userID)
		if err != nil {
			logger.Debug("Error getting steam id: " + err.Error())
			http.Error(w, "Error getting steam id", http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(steamId))
	}

}
