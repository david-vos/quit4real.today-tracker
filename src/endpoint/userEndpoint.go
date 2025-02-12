package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
	"quit4real.today/src/model"
	"strings"
)

type UserEndpoint struct {
	Router                     *mux.Router
	SteamService               *service.SteamService
	UserCommandHandler         *command.UserCommandHandler
	UserQueryHandler           *query.UserQueryHandler
	SubscriptionCommandHandler *command.SubscriptionCommandHandler
	AuthService                *service.AuthService
}

func (endpoint *UserEndpoint) User() {
	logger.Info("Trying to start the user endpoints")
	endpoint.Router.HandleFunc("/users", endpoint.RegisterHandler()).Methods("POST")
	endpoint.Router.HandleFunc("/users/${userName}/steamId", endpoint.GetSteamId()).Methods("GET")
	endpoint.Router.HandleFunc("/users/login", endpoint.LoginHandler()).Methods("POST")

	http.HandleFunc("/api/auth/steam", endpoint.LinkSteamAccountHandler())
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

		steamId, err := endpoint.SteamService.GetSteamIdFromVanityName(userID)
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
		errAddUser := endpoint.UserCommandHandler.Add(user)
		if errAddUser != nil {
			logger.Debug("Error adding user: " + errAddUser.Error())
			http.Error(w, "Error saving user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		return
	}
}

func (endpoint *UserEndpoint) LinkSteamAccountHandler() http.HandlerFunc {
	return endpoint.AuthService.AuthMiddleware(func(w http.ResponseWriter,
		r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}

		// Extract SteamID from OpenID response
		claimedID := r.FormValue("openid.claimed_id")
		const steamPrefix = "https://steamcommunity.com/openid/id/"
		if !strings.HasPrefix(claimedID, steamPrefix) {
			http.Error(w, "Invalid Steam response", http.StatusUnauthorized)
			return
		}
		steamID := strings.TrimPrefix(claimedID, steamPrefix)

		// Prepare verification request back to Steam
		params := r.URL.Query()
		params.Set("openid.mode", "check_authentication")

		resp, err := http.PostForm("https://steamcommunity.com/openid/login",
			params)
		if err != nil {
			http.Error(w, "Failed to verify OpenID",
				http.StatusInternalServerError)
			return
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logger.Fail("Failed to close response body")
				return
			}
		}(resp.Body)

		// Read entire response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, "Failed to read Steam response",
				http.StatusInternalServerError)
			return
		}
		if !strings.Contains(string(body), "is_valid:true") {
			http.Error(w, "Steam OpenID verification failed",
				http.StatusUnauthorized)
			return
		}

		// Extract username from the JWT
		tokenString := r.Header.Get("Authorization")
		username, err := endpoint.AuthService.GetFieldFromJWT(tokenString, "username")
		if err != nil {
			http.Error(w, "Error getting userId from user", http.StatusUnauthorized)
			return
		}

		// Retrieve user by username
		user, err := endpoint.UserQueryHandler.GetById(username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Update user with SteamID
		user.SteamID = steamID
		if err := endpoint.UserCommandHandler.Update(user); err != nil {
			logger.Debug("Error updating user with SteamID: " + err.Error())
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"steam_id": steamID})
	})
}
