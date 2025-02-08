package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"quit4real.today/logger"
	"quit4real.today/src/api"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
	"quit4real.today/src/model"
)

type UserEndpoint struct {
	Router                     *mux.Router
	SteamApi                   *api.SteamApi
	UserCommandHandler         *command.UserCommandHandler
	UserQueryHandler           *query.UserQueryHandler
	SubscriptionCommandHandler *command.SubscriptionCommandHandler
	AuthService                *service.AuthService
}

func (endpoint *UserEndpoint) User() {
	logger.Info("Trying to start the user endpoints")
	endpoint.Router.HandleFunc("/users", endpoint.RegisterHandler()).Methods("POST")
	endpoint.Router.HandleFunc("/users/${userName}/steamId", endpoint.GetSteamId()).Methods("GET")
	logger.Info("User endpoints started")
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

func (endpoint *UserEndpoint) LoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		var creds model.Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		user, err := endpoint.UserQueryHandler.GetById(creds.Username)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		}
		if !endpoint.AuthService.CheckPassword(user.Password, creds.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := endpoint.AuthService.GenerateJWT(creds.Username)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]string{"token": token})
		if err != nil {
			return
		}
	}
}

func (endpoint *UserEndpoint) RegisterHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		logger.Info("Got request to add a new user")

		var creds model.Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			http.Error(w, "Invalid JSON body", http.StatusBadRequest)
			return
		}

		if creds.Username == "" || creds.Password == "" {
			http.Error(w, "Missing username or password", http.StatusBadRequest)
			return
		}

		// Hash password before storing
		hashedPassword, err := endpoint.AuthService.HashPassword(creds.Password)
		if err != nil {
			http.Error(w, "Error hashing password", http.StatusInternalServerError)
			return
		}

		user := model.User{
			Name:     creds.Username,
			Password: string(hashedPassword),
		}
		var errAddUser = endpoint.UserCommandHandler.Add(user)
		if errAddUser != nil {
			logger.Debug("Error adding user: " + err.Error())
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}
