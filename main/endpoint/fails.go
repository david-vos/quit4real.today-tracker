package endpoint

import (
	"github.com/gorilla/mux"
	"project/logger"
)

type FailEndpoint struct {
	Router *mux.Router
}

func (f *FailEndpoint) Fail() {
	logger.Info("Trying to start fail endpoint")
	f.Router.HandleFunc("/fail/leaderboard").Methods("GET")
	logger.Info("Fail endpoint started")
}
