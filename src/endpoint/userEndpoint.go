package endpoint

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"quit4real.today/config"
	"quit4real.today/logger"
	"quit4real.today/src/handler/command"
	"quit4real.today/src/handler/query"
	"quit4real.today/src/handler/service"
	"quit4real.today/src/model"
)

type UserEndpoint struct {
	// Legacy code
	Router                     *mux.Router
	UserCommandHandler         *command.UserCommandHandlerImpl
	UserQueryHandler           *query.UserQueryHandlerImpl
	SubscriptionCommandHandler *command.SubscriptionCommandHandlerImpl
	// Services
	SteamService service.SteamService
	AuthService  service.AuthService
	UserService  service.UserService
}

func (endpoint *UserEndpoint) User() {
	logger.Info("Trying to start the user endpoints")
	endpoint.Router.HandleFunc("/users", endpoint.RegisterHandler()).Methods("POST")
	endpoint.Router.HandleFunc("/users/{userName}/steamId", endpoint.GetSteamId()).Methods("GET")
	endpoint.Router.HandleFunc("/users/login", endpoint.LoginHandler()).Methods("POST")

	endpoint.Router.HandleFunc("/api/auth/steam", endpoint.AuthService.AuthMiddleware(endpoint.SteamLoginHandler()))
	endpoint.Router.HandleFunc("/api/auth/steam/callback", endpoint.SteamCallbackHandler())
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
		err = json.NewEncoder(w).Encode(map[string]string{"steamId": steamId})
		if err != nil {
			http.Error(w, "Error getting steam id", http.StatusNoContent)
			return
		}

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
			return
		}
		if !endpoint.AuthService.CheckPassword(user.Password, creds.Password) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		token, err := endpoint.AuthService.GenerateJWT(user)
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]string{"token": token})
		if err != nil {
			http.Error(w, "Error generating token", http.StatusInternalServerError)
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

func (endpoint *UserEndpoint) SteamLoginHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		callbackURL := config.BackendUrl() + "/api/auth/steam/callback"

		openId := endpoint.AuthService.GetOpenId()
		redirectURL, err := openId.RedirectURL("https://steamcommunity.com/openid", callbackURL, "")
		if err != nil {
			http.Error(w, "OpenID Auth error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, redirectURL, http.StatusFound)
	}
}

func (endpoint *UserEndpoint) SteamCallbackHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fullURL := "https://" + r.Host + r.RequestURI
		logger.Info("Full URL: " + fullURL)

		openId := endpoint.AuthService.GetOpenId()
		id, err := openId.Verify(fullURL, nil, nil)
		if err != nil {
			logger.Fail("Failed to verify OpenID: " + err.Error())
			http.Error(w, "Failed to verify OpenID", http.StatusUnauthorized)
			return
		}

		// Extract SteamID from the OpenID URL
		steamID := id[len("https://steamcommunity.com/openid/id/"):]
		logger.Info("Steam ID verified: " + steamID)

		tokenString := r.Header.Get("Authorization")
		username, err := endpoint.AuthService.GetFieldFromJWT(tokenString, "username")
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		user, err := endpoint.UserQueryHandler.GetById(username)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		steamUserInfo, err := endpoint.SteamService.FetchUserInfo(steamID)
		if err != nil {
			http.Error(w, "Failed to get steam information", http.StatusNotFound)
			return
		}

		// update the user information to include steam info

		user.SteamUserName = steamUserInfo.PersonaName
		user.SteamID = steamID
		err = endpoint.UserCommandHandler.Update(user)
		if err != nil {
			http.Error(w, "Error updating user", http.StatusInternalServerError)
			return
		}

		// Add the trackers for steam to the user
		endpoint.UserService.CreateUserTrackers(user.SteamID)

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(map[string]string{"steamID": steamID})
		if err != nil {
			logger.Fail("Failed to encode JSON: " + err.Error()) // Log the error
			http.Error(w, "Failed to parse steamID", http.StatusUnauthorized)
			return
		}
	}
}
