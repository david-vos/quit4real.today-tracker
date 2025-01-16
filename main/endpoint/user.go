package endpoint

import (
	"github.com/gorilla/mux"
	"project/logger"
)

type UserEndpoint struct {
	Router *mux.Router
}

func (endpoint *UserEndpoint) User() {
	logger.Info("Trying to start the user endpoints")
	endpoint.Router.HandleFunc("/users").Methods("POST")
	endpoint.Router.HandleFunc("/user/{userID}/track/{gameID}").Methods("GET")
	endpoint.Router.HandleFunc("/user/{userID}/track/{gameID}").Methods("POST")
	logger.Info("User endpoints started")
}
